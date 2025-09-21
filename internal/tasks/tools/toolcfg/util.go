package toolcfg

import (
	"fmt"
	"os"
	"path/filepath"

	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// SetupConfigFile handles reading a Tool's config file from the embedded filesystem, and writing it
// to its target location.
func SetupConfigFile(t taskutil.Tool) error {
	cfgFileContents, err := Files.ReadFile(filepath.Base(t.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(t.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
