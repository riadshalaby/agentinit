# Plan

Status: **active**

Goal: replace the broken MCP server with a spawn-per-command session architecture that makes auto mode reliable, as scoped in `ROADMAP.md`.

## Scope

Rewrite `internal/mcp/` to use short-lived CLI subprocesses with provider-native session resume semantics. Add config loading, a disk-backed session store, provider adapters for Claude and Codex, and a new 7-tool MCP surface. Update scaffold templates and documentation.

## Acceptance Criteria

- `go build ./...` passes at every task boundary.
- `go test ./...` passes at every task boundary (coverage temporarily narrows at T-001, fully restored by T-006).
- `session_start` creates a named session, runs the provider CLI with the role prompt, persists metadata to `.ai/sessions.json`.
- `session_run` resumes the session, sends a command, blocks until the CLI process exits, returns full output in one call.
- Sessions survive MCP server restarts: `.ai/sessions.json` loaded on startup; any session marked `running` is reset to `errored`.
- Both adapters pass contract tests using Go test helper processes (no real CLI dependency in CI).
- In-process MCP client tests cover all 7 tools.
- E2E MCP handshake test passes.
- Manual mode scripts are unaffected.
- `agentinit init` scaffold produces updated PO prompt with new tool names.

## Key Design Decisions (resolved)

| Decision | Choice |
|----------|--------|
| Execution model | Spawn-per-command — no long-lived pipes, no PTY |
| `session_run` | Synchronous — blocks until CLI process exits |
| MCP-startable roles | `implement` and `review` only (validated in manager) |
| Session persistence | `.ai/sessions.json`, gitignored |
| Tool name migration | Clean break — old 5-tool surface removed entirely |
| Naming conflict resolution | Delete old session files in T-001, define clean types immediately |

## Sequencing Rationale

`session.go` exports `Session`, `SessionStatus`, `SessionManager`, and `SessionInfo` — all names needed by the new design. Rather than staging with temporary names, T-001 deletes the old files and stubs the server to compile, giving all subsequent tasks a clean base.

## Implementation Phases

---

### T-001 — Foundation: delete legacy session implementation, define domain types

**Files changed:**
- DELETE `internal/mcp/session.go`
- DELETE `internal/mcp/session_test.go`
- CREATE `internal/mcp/types.go`
- MODIFY `internal/mcp/server.go`
- MODIFY `internal/mcp/tools.go`
- MODIFY `internal/mcp/server_test.go`

**`types.go` — define the domain model:**

```go
package mcp

import "time"

type SessionStatus string

const (
    StatusIdle    SessionStatus = "idle"
    StatusRunning SessionStatus = "running"
    StatusErrored SessionStatus = "errored"
    StatusStopped SessionStatus = "stopped"
)

type ProviderState struct {
    SessionID string `json:"session_id,omitempty"`
}

type Session struct {
    Name          string        `json:"name"`
    Role          string        `json:"role"`
    Provider      string        `json:"provider"`
    Model         string        `json:"model,omitempty"`
    Status        SessionStatus `json:"status"`
    ProviderState ProviderState `json:"provider_state"`
    CreatedAt     time.Time     `json:"created_at"`
    LastActiveAt  time.Time     `json:"last_active_at"`
    RunCount      int           `json:"run_count"`
    Error         string        `json:"error,omitempty"`
}

type SessionInfo struct {
    Name     string        `json:"name"`
    Role     string        `json:"role"`
    Provider string        `json:"provider"`
    Status   SessionStatus `json:"status"`
    RunCount int           `json:"run_count"`
    Error    string        `json:"error,omitempty"`
}

func (s *Session) info() SessionInfo {
    return SessionInfo{
        Name:     s.Name,
        Role:     s.Role,
        Provider: s.Provider,
        Status:   s.Status,
        RunCount: s.RunCount,
        Error:    s.Error,
    }
}
```

**`server.go` — remove all references to the deleted `SessionManager`; stub a placeholder:**

```go
package mcp

import (
    "context"
    "log/slog"

    mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
    server *mcpserver.MCPServer
    logger *slog.Logger
}

func NewServer(version string) *Server {
    logger, err := NewFileLogger(defaultMCPLogPath)
    if err != nil {
        logger = newDiscardLogger()
    }
    return newServer(version, logger)
}

func newServer(version string, logger *slog.Logger) *Server {
    if logger == nil {
        logger = newDiscardLogger()
    }
    srv := mcpserver.NewMCPServer(serverName, version, mcpserver.WithToolCapabilities(false))
    registerTools(srv, logger)
    return &Server{server: srv, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
    _ = ctx
    return serveStdio(s.server)
}
```

**`tools.go` — register all 7 stub tools (return `"not implemented"` errors):**

Register these 7 tools with correct names and descriptions but stub handlers:
- `session_start` (args: `name` string required, `role` string required, `provider` string optional)
- `session_run` (args: `name` string required, `command` string required, `timeout_seconds` number optional)
- `session_status` (args: `name` string required)
- `session_list` (no args)
- `session_stop` (args: `name` string required)
- `session_reset` (args: `name` string required)
- `session_delete` (args: `name` string required)

All stub handlers return `mcpproto.NewToolResultErrorf("not implemented")`.

Retain the `jsonResult` helper function.

**`server_test.go` — update to reflect stub state:**

- Keep `TestNewServerRespondsToInitialize`: unchanged.
- Update `TestNewServerRegistersSessionTools`: change tool count assertion from `5` to `7`; verify log file is created.
- Remove `TestServerSessionToolsLifecycle`: it depended on old `SessionManager` and will be fully replaced in T-006.
- Remove all test helpers (`testLauncher`, `testSpawnLauncher`, `testLogger`, `TestHelperSessionProcess`, `TestHelperSpawnProcess`, `assertStructuredToolResult`, `containsAll`) — they will be reintroduced in T-004 and T-006 with new signatures.

**Acceptance criteria:**
- `go build ./...` passes.
- `go test ./...` passes.
- `internal/mcp` package has no reference to old `Session`, `SpawnSession`, `launcherFunc`, `spawnLauncherFunc`, or `spawnRequest` types.
- `server_test.go` confirms 7 tools registered.

---

### T-002 — Config layer

**Files changed:**
- CREATE `internal/mcp/config.go`
- CREATE `internal/mcp/config_test.go`

**`config.go`:**

```go
package mcp

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
)

var validProviders = map[string]struct{}{"claude": {}, "codex": {}}
var validRoles    = map[string]struct{}{"implement": {}, "review": {}}

type Config struct {
    Roles    map[string]RoleConfig `json:"roles"`
    Defaults ProviderDefaults      `json:"defaults,omitempty"`
}

type RoleConfig struct {
    Provider string `json:"agent,omitempty"`
    Model    string `json:"model,omitempty"`
    Effort   string `json:"effort,omitempty"`
}

type ProviderDefaults struct {
    Claude ClaudeDefaults `json:"claude,omitempty"`
    Codex  CodexDefaults  `json:"codex,omitempty"`
}

type ClaudeDefaults struct {
    PermissionMode string `json:"permission_mode,omitempty"`
}

type CodexDefaults struct {
    Sandbox       string `json:"sandbox,omitempty"`
    NetworkAccess bool   `json:"network_access,omitempty"`
}

// LoadConfig reads .ai/config.json from cwd. A missing file is not an error;
// it returns a zero-value Config. An invalid file returns an error.
func LoadConfig(cwd string) (Config, error) {
    path := filepath.Join(cwd, ".ai", "config.json")
    data, err := os.ReadFile(path)
    if errors.Is(err, os.ErrNotExist) {
        return Config{}, nil
    }
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, fmt.Errorf("parse config: %w", err)
    }
    return cfg, nil
}

// ProviderForRole returns the configured provider for a role, defaulting to "claude".
func (c Config) ProviderForRole(role string) string {
    if rc, ok := c.Roles[role]; ok && rc.Provider != "" {
        return rc.Provider
    }
    return "claude"
}

// ModelForRole returns the configured model for a role. Empty string means
// the provider's own default.
func (c Config) ModelForRole(role string) string {
    if rc, ok := c.Roles[role]; ok {
        return rc.Model
    }
    return ""
}

// EffortForRole returns the configured effort for a role (Claude-specific).
func (c Config) EffortForRole(role string) string {
    if rc, ok := c.Roles[role]; ok {
        return rc.Effort
    }
    return ""
}
```

Note: `validProviders` and `validRoles` are package-level vars reused by both the config layer and the session manager.

**`config_test.go`** — cover:
- `LoadConfig` from a temp dir with no `.ai/config.json` → zero-value Config, no error.
- `LoadConfig` with the existing project template JSON → correct provider/model/effort per role.
- `LoadConfig` with malformed JSON → error.
- `ProviderForRole` for a known role → returns configured value.
- `ProviderForRole` for an unknown role → returns `"claude"`.
- `ModelForRole` for a role with model → returns model; without → returns `""`.
- `EffortForRole` for a role with effort → returns effort; without → returns `""`.
- `Defaults` block: config with `defaults.claude.permission_mode` set → field is accessible.

**Acceptance criteria:**
- `go test ./internal/mcp/... -run TestConfig` passes.
- `go test ./...` passes.

---

### T-003 — Session store

**Files changed:**
- CREATE `internal/mcp/store.go`
- CREATE `internal/mcp/store_test.go`

**`store.go`:**

```go
package mcp

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

const defaultSessionsPath = ".ai/sessions.json"

// Store persists session metadata to a JSON file.
// All methods are safe for concurrent use.
type Store struct {
    path string
    mu   sync.Mutex
}

func NewStore(path string) *Store {
    if path == "" {
        path = defaultSessionsPath
    }
    return &Store{path: path}
}

// Load reads the store file and returns all sessions.
// A missing file returns an empty map without error.
// A corrupt file returns an error.
func (s *Store) Load() (map[string]*Session, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.load()
}

// Put writes or replaces a session by name.
func (s *Store) Put(session *Session) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    sessions, err := s.load()
    if err != nil {
        return err
    }
    sessions[session.Name] = session
    return s.save(sessions)
}

// Get returns a session by name, or an error if not found.
func (s *Store) Get(name string) (*Session, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    sessions, err := s.load()
    if err != nil {
        return nil, err
    }
    session, ok := sessions[name]
    if !ok {
        return nil, fmt.Errorf("session %q not found", name)
    }
    return session, nil
}

// Delete removes a session by name. No-op if not found.
func (s *Store) Delete(name string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    sessions, err := s.load()
    if err != nil {
        return err
    }
    delete(sessions, name)
    return s.save(sessions)
}

// List returns all sessions in undefined order.
func (s *Store) List() ([]*Session, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    sessions, err := s.load()
    if err != nil {
        return nil, err
    }
    out := make([]*Session, 0, len(sessions))
    for _, s := range sessions {
        out = append(out, s)
    }
    return out, nil
}

func (s *Store) load() (map[string]*Session, error) {
    data, err := os.ReadFile(s.path)
    if errors.Is(err, os.ErrNotExist) {
        return make(map[string]*Session), nil
    }
    if err != nil {
        return nil, fmt.Errorf("read sessions file: %w", err)
    }
    var sessions map[string]*Session
    if err := json.Unmarshal(data, &sessions); err != nil {
        return nil, fmt.Errorf("parse sessions file: %w", err)
    }
    if sessions == nil {
        sessions = make(map[string]*Session)
    }
    return sessions, nil
}

func (s *Store) save(sessions map[string]*Session) error {
    if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
        return fmt.Errorf("create sessions directory: %w", err)
    }
    data, err := json.MarshalIndent(sessions, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal sessions: %w", err)
    }
    return os.WriteFile(s.path, data, 0o644)
}
```

**`store_test.go`** — cover:
- `Load` on a missing file → empty map, no error.
- `Put` then `Get` → same session returned.
- `Put` then `List` → session appears in list.
- `Put` then `Delete` then `Get` → error.
- `Put` with two sessions → `List` returns both.
- Corrupt JSON → `Load` returns error.
- `Put` creates parent directory if missing (temp path like `/tmp/x/y/sessions.json`).

**Acceptance criteria:**
- `go test ./internal/mcp/... -run TestStore` passes.
- `go test ./...` passes.

---

### T-004 — Provider adapters

**Files changed:**
- CREATE `internal/mcp/adapter.go`
- CREATE `internal/mcp/adapter_codex.go`
- CREATE `internal/mcp/adapter_claude.go`
- CREATE `internal/mcp/adapter_test.go`

**`adapter.go` — interface and options:**

```go
package mcp

import (
    "context"
    "time"
)

// StartOpts are passed to Adapter.Start.
type StartOpts struct {
    PromptFile string        // path to role prompt file (required)
    Model      string        // provider-specific model string (optional)
    Effort     string        // claude: --effort value (optional)
    Timeout    time.Duration // 0 means no timeout
}

// RunOpts are passed to Adapter.Run.
type RunOpts struct {
    Model   string
    Timeout time.Duration // 0 means no timeout
}

// Adapter handles provider-specific CLI invocation.
// Each method spawns a short-lived subprocess and returns its full output.
type Adapter interface {
    // Start runs the initial CLI invocation with the role system prompt.
    // It updates session.ProviderState in place.
    Start(ctx context.Context, session *Session, opts StartOpts) (output string, err error)

    // Run resumes the session with a command.
    // It updates session.ProviderState in place.
    Run(ctx context.Context, session *Session, command string, opts RunOpts) (output string, err error)

    // Stop kills the process identified by session.ProviderState if it is
    // currently running. No-op if nothing is running.
    Stop(ctx context.Context, session *Session) error
}
```

**`adapter_codex.go`:**

```go
package mcp

import (
    "context"
    "fmt"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

var codexSessionIDPattern = regexp.MustCompile(`(?m)^session id:\s+(\S+)$`)

type CodexAdapter struct {
    cwd     string
    sandbox string
    network bool
    exec    func(ctx context.Context, args []string, stdin string) (string, error)
}

func NewCodexAdapter(cwd string, defaults CodexDefaults) *CodexAdapter {
    sandbox := defaults.Sandbox
    if sandbox == "" {
        sandbox = "workspace-write"
    }
    a := &CodexAdapter{cwd: cwd, sandbox: sandbox, network: defaults.NetworkAccess}
    a.exec = a.defaultExec
    return a
}

func (a *CodexAdapter) Start(ctx context.Context, session *Session, opts StartOpts) (string, error) {
    prompt, err := readPromptFile(opts.PromptFile)
    if err != nil {
        return "", err
    }
    args := []string{"exec", "--sandbox", a.sandbox}
    if a.network {
        args = append(args, "-c", fmt.Sprintf("sandbox_%s.network_access=true", strings.ReplaceAll(a.sandbox, "-", "_")))
    }
    if opts.Model != "" {
        args = append(args, "-m", opts.Model)
    }
    args = append(args, "-")
    output, err := a.exec(ctx, args, prompt)
    if sessionID := extractCodexSessionID(output); sessionID != "" {
        session.ProviderState.SessionID = sessionID
    }
    return output, err
}

func (a *CodexAdapter) Run(ctx context.Context, session *Session, command string, opts RunOpts) (string, error) {
    if session.ProviderState.SessionID == "" {
        return "", fmt.Errorf("session %q has no provider session ID; call Start first", session.Name)
    }
    args := []string{"exec", "resume", session.ProviderState.SessionID}
    if a.network {
        args = append(args, "-c", fmt.Sprintf("sandbox_%s.network_access=true", strings.ReplaceAll(a.sandbox, "-", "_")))
    }
    if opts.Model != "" {
        args = append(args, "-m", opts.Model)
    }
    args = append(args, "-")
    output, err := a.exec(ctx, args, command)
    if sessionID := extractCodexSessionID(output); sessionID != "" {
        session.ProviderState.SessionID = sessionID
    }
    return output, err
}

func (a *CodexAdapter) Stop(_ context.Context, _ *Session) error {
    // Spawn-per-command: no persistent process to stop.
    return nil
}

func (a *CodexAdapter) defaultExec(ctx context.Context, args []string, stdin string) (string, error) {
    cmd := exec.CommandContext(ctx, "codex", args...)
    cmd.Dir = a.cwd
    cmd.Stdin = strings.NewReader(stdin + "\n")
    out, err := cmd.CombinedOutput()
    return string(out), err
}

func extractCodexSessionID(output string) string {
    matches := codexSessionIDPattern.FindStringSubmatch(output)
    if len(matches) != 2 {
        return ""
    }
    return matches[1]
}
```

**`adapter_claude.go`:**

```go
package mcp

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
)

type ClaudeAdapter struct {
    cwd            string
    permissionMode string
    exec           func(ctx context.Context, args []string) (string, error)
}

func NewClaudeAdapter(cwd string, defaults ClaudeDefaults) *ClaudeAdapter {
    pm := defaults.PermissionMode
    if pm == "" {
        pm = "acceptEdits"
    }
    a := &ClaudeAdapter{cwd: cwd, permissionMode: pm}
    a.exec = a.defaultExec
    return a
}

func (a *ClaudeAdapter) Start(ctx context.Context, session *Session, opts StartOpts) (string, error) {
    if session.ProviderState.SessionID == "" {
        return "", fmt.Errorf("session %q has no session ID; caller must set one before Start", session.Name)
    }
    args := []string{"-p",
        "--session-id", session.ProviderState.SessionID,
        "--permission-mode", a.permissionMode,
    }
    if opts.PromptFile != "" {
        args = append(args, "--system-prompt-file", opts.PromptFile)
    }
    if opts.Model != "" {
        args = append(args, "--model", opts.Model)
    }
    if opts.Effort != "" {
        args = append(args, "--effort", opts.Effort)
    }
    args = append(args, "You are now in WAIT_FOR_USER_START state.")
    return a.exec(ctx, args)
}

func (a *ClaudeAdapter) Run(ctx context.Context, session *Session, command string, opts RunOpts) (string, error) {
    if session.ProviderState.SessionID == "" {
        return "", fmt.Errorf("session %q has no provider session ID; call Start first", session.Name)
    }
    args := []string{"-p",
        "--session-id", session.ProviderState.SessionID,
        "--permission-mode", a.permissionMode,
    }
    if opts.Model != "" {
        args = append(args, "--model", opts.Model)
    }
    args = append(args, command)
    return a.exec(ctx, args)
}

func (a *ClaudeAdapter) Stop(_ context.Context, _ *Session) error {
    // Spawn-per-command: no persistent process to stop.
    return nil
}

func (a *ClaudeAdapter) defaultExec(ctx context.Context, args []string) (string, error) {
    cmd := exec.CommandContext(ctx, "claude", args...)
    cmd.Dir = a.cwd
    out, err := cmd.CombinedOutput()
    return string(out), err
}
```

Note on `session.ProviderState.SessionID` for the Claude adapter: the manager must assign a UUID to `session.ProviderState.SessionID` before calling `Start`. This is the manager's responsibility, not the adapter's. The adapter uses the ID it receives.

**`adapter_test.go`** — use Go test helper processes (same pattern as the deleted `TestHelperSpawnProcess`):

Define two test helpers:
- `TestHelperCodexProcess`: activated by `GO_WANT_HELPER_CODEX=1`. Reads args to determine `start` vs `resume`. On start, prints `session id: test-session-abc`. On resume, prints the stdin content as the response. Exits 0.
- `TestHelperClaudeProcess`: activated by `GO_WANT_HELPER_CLAUDE=1`. Reads args and prints `--session-id` value + command back as response. Exits 0.

Wire the adapters to use these helpers via the `exec` field (same technique as `testSpawnLauncher` in the deleted tests).

Contract tests for each adapter:
- `TestCodexAdapterStart`: captures session ID from test helper output; updates `session.ProviderState.SessionID`.
- `TestCodexAdapterRun`: passes correct `resume` args; passes command via stdin; captures output.
- `TestCodexAdapterRunNoSessionID`: returns error when session has no session ID.
- `TestClaudeAdapterStart`: passes `--session-id` and `--system-prompt-file` args.
- `TestClaudeAdapterRun`: passes `--session-id` from existing provider state.
- `TestClaudeAdapterRunNoSessionID`: returns error when session has no session ID.

**Acceptance criteria:**
- `go test ./internal/mcp/... -run TestAdapter` passes.
- `go test ./internal/mcp/... -run TestHelper` passes (helpers run only when env var is set).
- `go test ./...` passes.

---

### T-005 — Session manager

**Files changed:**
- CREATE `internal/mcp/manager.go`
- CREATE `internal/mcp/manager_test.go`

**`manager.go`:**

```go
package mcp

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/google/uuid"
)

// SessionManager owns the named session registry.
// It is the single entry point for all session lifecycle operations.
type SessionManager struct {
    store    *Store
    adapters map[string]Adapter
    config   Config
    cwd      string
    mu       sync.Mutex
    running  map[string]context.CancelFunc // name -> cancel for in-flight Run
    logger   *slog.Logger
}

func NewSessionManager(store *Store, adapters map[string]Adapter, config Config, cwd string, logger *slog.Logger) *SessionManager {
    if logger == nil {
        logger = newDiscardLogger()
    }
    m := &SessionManager{
        store:    store,
        adapters: adapters,
        config:   config,
        cwd:      cwd,
        running:  make(map[string]context.CancelFunc),
        logger:   logger,
    }
    m.recoverStaleRunning()
    return m
}

// recoverStaleRunning marks any session persisted as StatusRunning as StatusErrored
// on startup (the previous MCP server process died mid-run).
func (m *SessionManager) recoverStaleRunning() { ... }

// StartSession creates a new named session, runs the provider CLI with the role
// prompt, and persists the session. Returns an error if the name is already in use.
func (m *SessionManager) StartSession(ctx context.Context, name, role, provider string) (SessionInfo, string, error) {
    if err := validateRole(role); err != nil {
        return SessionInfo{}, "", err
    }
    if err := validateProvider(provider); err != nil {
        return SessionInfo{}, "", err
    }
    // Fail if session name already exists
    // Build Session, set ProviderState.SessionID to new UUID (for claude)
    // Call adapter.Start, update session.Status=idle, persist
    ...
}

// RunSession sends a command to an existing session synchronously.
// Returns an error if the session is already running.
func (m *SessionManager) RunSession(ctx context.Context, name, command string, timeout time.Duration) (SessionInfo, string, error) {
    // Acquire run lock for name
    // Set status=running, persist
    // Call adapter.Run, update status=idle or errored
    // Persist, release lock
    ...
}

// StopSession cancels any in-flight RunSession for the named session.
func (m *SessionManager) StopSession(name string) (SessionInfo, error) { ... }

// ResetSession clears provider state so the next Run starts a fresh conversation.
func (m *SessionManager) ResetSession(name string) (SessionInfo, error) { ... }

// DeleteSession removes the session entirely.
func (m *SessionManager) DeleteSession(name string) error { ... }

// GetSession returns the current SessionInfo for a named session.
func (m *SessionManager) GetSession(name string) (SessionInfo, error) { ... }

// ListSessions returns info for all tracked sessions.
func (m *SessionManager) ListSessions() ([]SessionInfo, error) { ... }

func validateProvider(provider string) error {
    if _, ok := validProviders[provider]; !ok {
        return fmt.Errorf("unsupported provider %q", provider)
    }
    return nil
}

func validateRole(role string) error {
    if _, ok := validRoles[role]; !ok {
        return fmt.Errorf("unsupported role %q: must be one of: implement, review", role)
    }
    return nil
}
```

**Prompt file resolution** — add a package-level helper (moved here from deleted `session.go`):

```go
func promptFileForRole(cwd, role string) (string, error) {
    var name string
    switch role {
    case "implement":
        name = "implementer.md"
    case "review":
        name = "reviewer.md"
    default:
        return "", fmt.Errorf("no prompt file for role %q", role)
    }
    path := filepath.Join(cwd, ".ai", "prompts", name)
    if _, err := os.Stat(path); err != nil {
        return "", fmt.Errorf("locate prompt file for role %q: %w", role, err)
    }
    return path, nil
}

func readPromptFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("read prompt file %q: %w", path, err)
    }
    return string(data), nil
}
```

**Run concurrency:** `m.running` holds a `context.CancelFunc` per session name. If `RunSession` is called while that name has an entry in `m.running`, return an error: `"session %q is already running"`. The cancel func is called by `StopSession`.

**`manager_test.go`** — use test adapters:

Implement a `testAdapter` struct that:
- `Start`: sets `session.ProviderState.SessionID = "test-session-id"`, returns `"WAIT_FOR_USER_START"`, nil.
- `Run`: returns `fmt.Sprintf("response: %s", command)`, nil.
- `Stop`: no-op.

Tests:
- `TestManagerStartSession`: creates session, name present in `List`.
- `TestManagerStartDuplicateName`: second `StartSession` with same name returns error.
- `TestManagerRunSession`: `Run` returns expected output; `RunCount` increments; `LastActiveAt` updated.
- `TestManagerRunConcurrent`: second `RunSession` while first is in-flight returns error.
- `TestManagerStopSession`: `StopSession` cancels the in-flight run context.
- `TestManagerResetSession`: clears `ProviderState.SessionID`; `Status` → `idle`.
- `TestManagerDeleteSession`: session absent from `List` and `Get` after delete.
- `TestManagerRestartRecovery`: create `sessions.json` with a session at `StatusRunning`; construct manager; that session is `StatusErrored` after construction.
- `TestManagerStartInvalidRole`: error for role not in `validRoles`.
- `TestManagerStartInvalidProvider`: error for provider not in `validProviders`.

**Acceptance criteria:**
- `go test ./internal/mcp/... -run TestManager` passes.
- `go test ./...` passes.

---

### T-006 — MCP tool surface

**Files changed:**
- REWRITE `internal/mcp/tools.go`
- MODIFY `internal/mcp/server.go` (wire real manager; update `newServer` signature)
- REWRITE `internal/mcp/server_test.go`
- MODIFY `e2e/e2e_test.go` (update tool count assertion)

**`tools.go`** — replace stub handlers with real implementations:

```
session_start:
  args: name (string, required), role (string, required), provider (string, optional — defaults from config for the role)
  calls: manager.StartSession
  on success: jsonResult with SessionInfo + initial output text

session_run:
  args: name (string, required), command (string, required), timeout_seconds (number, optional, default 300)
  calls: manager.RunSession
  on success: jsonResult with SessionInfo + output text

session_status:
  args: name (string, required)
  calls: manager.GetSession
  on success: jsonResult with SessionInfo

session_list:
  args: (none)
  calls: manager.ListSessions
  on success: jsonResult with []SessionInfo

session_stop:
  args: name (string, required)
  calls: manager.StopSession
  on success: jsonResult with SessionInfo

session_reset:
  args: name (string, required)
  calls: manager.ResetSession
  on success: jsonResult with SessionInfo

session_delete:
  args: name (string, required)
  calls: manager.DeleteSession
  on success: jsonResult with {"name": name, "deleted": true}
```

Provider defaulting in `session_start`: if `provider` arg is empty or absent, call `config.ProviderForRole(role)` to fill it.

**`server.go`** — update `NewServer` and `newServer`:

```go
func NewServer(version string) *Server {
    logger, err := NewFileLogger(defaultMCPLogPath)
    if err != nil {
        logger = newDiscardLogger()
    }

    cwd, _ := os.Getwd()
    cfg, _ := LoadConfig(cwd)

    store := NewStore(filepath.Join(cwd, defaultSessionsPath))
    adapters := map[string]Adapter{
        "claude": NewClaudeAdapter(cwd, cfg.Defaults.Claude),
        "codex":  NewCodexAdapter(cwd, cfg.Defaults.Codex),
    }
    manager := NewSessionManager(store, adapters, cfg, cwd, logger)

    return newServer(version, manager, cfg, logger)
}

func newServer(version string, manager *SessionManager, cfg Config, logger *slog.Logger) *Server {
    ...
    registerTools(srv, manager, cfg, logger)
    ...
}
```

**`server_test.go`** — reintroduce full lifecycle test using `testAdapter`:

- `TestNewServerRespondsToInitialize`: unchanged.
- `TestNewServerRegistersSessionTools`: tool count = 7; log file created.
- `TestServerSessionToolsLifecycle`: in-process MCP client calls:
  1. `session_start(name="implementer", role="implement", provider="codex")` → not error, `session_id` in result
  2. `session_run(name="implementer", command="next_task T-001")` → not error, output contains response
  3. `session_status(name="implementer")` → status is "idle"
  4. `session_list()` → contains "implementer"
  5. `session_start(name="implementer", role="implement", provider="codex")` duplicate → IsError
  6. `session_reset(name="implementer")` → not error
  7. `session_delete(name="implementer")` → not error
  8. `session_status(name="implementer")` → IsError (not found)

Use `newServer` with a pre-built `SessionManager` using `testAdapter`. Use a temp dir so the store writes to a temp path.

Reintroduce `jsonResult`, `assertStructuredToolResult`, `containsAll` helpers.

**`e2e/e2e_test.go`** — `TestMCPInitializeHandshake` is unchanged. No tool-count assertion exists in e2e; no changes needed beyond confirming the test still passes.

**Acceptance criteria:**
- `go test ./internal/mcp/...` passes fully.
- `go test -tags e2e ./e2e/...` passes.
- `go test ./...` passes.
- In-process MCP client test covers all 7 tools.

---

### T-007 — Template and documentation updates

**Files changed:**
- MODIFY `internal/template/templates/base/ai/prompts/po.md.tmpl`
- MODIFY `internal/template/templates/base/ai/config.json.tmpl`
- MODIFY `internal/template/templates/base/gitignore.tmpl`
- MODIFY `.ai/prompts/po.md` (this repo's own PO prompt)
- MODIFY `.ai/config.json` (this repo's own config — add `defaults` block example)
- MODIFY `.gitignore` (add `.ai/sessions.json`)
- MODIFY `README.md` (MCP tools table, config schema section, migration note)

**PO prompt changes** (both the template and the live file):

Replace the MCP tools section. Old tools (`start_session`, `send_command`, `get_output`, `list_sessions`, `stop_session`) become:

```
Use the agentinit MCP server tools to coordinate the other role sessions:
  - `session_start`  — create and initialize a named session
  - `session_run`    — send a command and receive the full output (synchronous)
  - `session_status` — check the current status of a session
  - `session_list`   — list all tracked sessions
  - `session_stop`   — cancel an in-flight run
  - `session_reset`  — clear provider state so next run starts a fresh conversation
  - `session_delete` — remove a session entirely
```

Replace the interaction pattern. The `send_command` + polling `get_output` loop is removed. The new pattern:

```
1. Re-read `.ai/TASKS.md`.
2. Decide the next deterministic action from the board state.
3. Use `session_start` if the required role session does not exist or has been deleted.
4. Use `session_run(name, command)` to send the role command and receive the full output.
5. Re-read `.ai/TASKS.md` to confirm the status transition before deciding the next command.
```

Session naming convention for PO prompt: use `"implementer"` for the implement session and `"reviewer"` for the review session.

Session start example for PO prompt:
```
session_start(name="implementer", role="implement")  // provider defaults from .ai/config.json
session_start(name="reviewer",    role="review")
```

**`config.json.tmpl`** — add `defaults` block:

```json
{
  "roles": {
    "plan": {
      "agent": "claude",
      "model": "sonnet",
      "effort": "medium"
    },
    "implement": {
      "agent": "codex",
      "model": "gpt-5.4"
    },
    "review": {
      "agent": "claude",
      "model": "sonnet",
      "effort": "medium"
    }
  },
  "defaults": {
    "claude": {
      "permission_mode": "acceptEdits"
    },
    "codex": {
      "sandbox": "workspace-write",
      "network_access": true
    }
  }
}
```

**`gitignore.tmpl`** — add `.ai/sessions.json` to the gitignore template.

**`.gitignore`** — add `.ai/sessions.json` to this repo's gitignore.

**`README.md`** — update the MCP section:
- Tools table: replace 5-tool table with 7-tool table.
- Add a brief note that `session_run` is synchronous (no polling required).
- Update the `config.json` schema description to mention the `defaults` block.
- Add a migration note at the end of the MCP section:

  > **0.7.0 Migration:** The MCP tool surface has been renamed and consolidated. Run `agentinit update` to get the updated PO prompt. `session_run` replaces the old `send_command` + `get_output` polling loop. Sessions are now named and persist across restarts in `.ai/sessions.json`.

**Acceptance criteria:**
- `agentinit init demo --no-git --dir /tmp/test-scaffold` produces:
  - PO prompt containing `session_start`, `session_run`, `session_delete`.
  - `.gitignore` containing `.ai/sessions.json`.
  - `config.json` containing `"defaults"` block.
- `agentinit update --dry-run` in an existing project reports PO prompt would update.
- README MCP tools table lists all 7 new tool names.
- `go test ./...` passes.

---

## Validation

Run after every task:
```
go fmt ./...
go vet ./...
go test ./...
```

Run before marking the cycle done:
```
go test -tags e2e ./e2e/...
```
