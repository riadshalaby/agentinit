# ROADMAP

Goal: consolidate manual and auto workflows into a single, unified workflow model.

## Priority 1 — Unified Scaffold

Objective: every scaffolded project gets all five roles (PO, planner, implementer, reviewer, tester) regardless of how the user interacts. The `--workflow manual|auto` flag is removed; the distinction becomes a runtime choice, not a scaffold-time decision.

### What changes

1. **Remove the `--workflow` flag from `agentinit init`.**
   - The wizard no longer asks manual vs auto.
   - Every project gets the full set of scripts and prompts including `ai-po.sh` and `po.md`.

2. **Always scaffold the PO role.**
   - `scripts/ai-po.sh`, `.ai/prompts/po.md`, and the MCP config are generated for every project.
   - The PO script is simply unused when the user drives manually.

3. **Two runtime modes, same scaffold.**
   - **Manual**: user opens separate terminals for planner, implementer, reviewer, tester and drives the cycle with session commands.
   - **Auto**: user runs `scripts/ai-po.sh` which starts the MCP server and a PO agent that orchestrates the other roles.

4. **Update templates and overlays.**
   - Remove any conditional logic gated on workflow type in scaffold, template, and overlay code.
   - Ensure base templates include PO artifacts unconditionally.

5. **Update documentation.**
   - `README.md.tmpl`: describe both modes as runtime options under a single setup.
   - `.ai/AGENTS.md.tmpl`: remove workflow-type branching; document both modes side by side.
   - CLI help text: reflect the simplified init interface.

