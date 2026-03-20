package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteFiles writes all rendered template content to disk under targetDir.
// Shell scripts get chmod +x.
func WriteFiles(targetDir string, files map[string]string) error {
	for relPath, content := range files {
		absPath := filepath.Join(targetDir, relPath)

		dir := filepath.Dir(absPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}

		perm := os.FileMode(0o644)
		if strings.HasSuffix(relPath, ".sh") {
			perm = 0o755
		}

		if err := os.WriteFile(absPath, []byte(content), perm); err != nil {
			return fmt.Errorf("write %s: %w", absPath, err)
		}
	}
	return nil
}
