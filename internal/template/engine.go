package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// filenameMapping translates template filenames to output filenames.
var filenameMapping = map[string]string{
	"gitignore.tmpl":     ".gitignore",
	"gitattributes.tmpl": ".gitattributes",
}

// RenderAll renders all base templates with the given data.
// It first renders any overlay fragments and injects results into data,
// then renders base templates and returns a map of relative output path -> content.
func RenderAll(data *ProjectData) (map[string]string, error) {
	result := make(map[string]string)

	// Step 1: Render overlay gitignore fragment if it exists.
	if data.ProjectType != "" {
		fragPath := fmt.Sprintf("templates/overlays/%s/gitignore_extra.tmpl", data.ProjectType)
		content, err := fs.ReadFile(TemplateFS, fragPath)
		if err == nil {
			rendered, err := renderTemplate(string(content), fragPath, data)
			if err != nil {
				return nil, fmt.Errorf("render overlay fragment %s: %w", fragPath, err)
			}
			data.GitignoreExtra = rendered
		}
	}

	// Step 2: Walk and render base templates.
	err := fs.WalkDir(TemplateFS, "templates/base", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(TemplateFS, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		rendered, err := renderTemplate(string(content), path, data)
		if err != nil {
			return fmt.Errorf("render %s: %w", path, err)
		}
		if strings.TrimSpace(rendered) == "" {
			return nil
		}

		// Convert template path to output path.
		outPath := strings.TrimPrefix(path, "templates/base/")
		outPath = strings.TrimSuffix(outPath, ".tmpl")

		// Map ai/ and claude/ directory prefixes to dot-directories in output.
		if strings.HasPrefix(outPath, "ai/") {
			outPath = ".ai/" + strings.TrimPrefix(outPath, "ai/")
		}
		if strings.HasPrefix(outPath, "claude/") {
			outPath = ".claude/" + strings.TrimPrefix(outPath, "claude/")
		}

		// Apply filename mappings for dotfiles.
		base := filepath.Base(path)
		if mapped, ok := filenameMapping[base]; ok {
			outPath = filepath.Join(filepath.Dir(outPath), mapped)
		}

		result[outPath] = rendered
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func renderTemplate(tmplContent, name string, data any) (string, error) {
	funcMap := template.FuncMap{
		"indent": func(spaces int, s string) string {
			pad := strings.Repeat(" ", spaces)
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = pad + line
				}
			}
			return strings.Join(lines, "\n")
		},
		"permissionRules": func(data *ProjectData) string {
			seen := make(map[string]struct{})
			commands := make([]string, 0, len(data.ToolPermissions)+len(data.ValidationCommands)+2)

			add := func(command string) {
				if _, ok := seen[command]; ok {
					return
				}
				seen[command] = struct{}{}
				commands = append(commands, command)
			}

			for _, command := range data.ToolPermissions {
				add(command)
			}
			for _, validation := range data.ValidationCommands {
				add(validation.Command)
			}
			add("git add")
			add("git commit")
			rules := make([]string, 0, len(commands))
			for _, command := range commands {
				rules = append(rules, fmt.Sprintf("\"Bash(%s:*)\"", command))
			}
			return strings.Join(rules, ",\n      ")
		},
	}

	t, err := template.New(name).Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
