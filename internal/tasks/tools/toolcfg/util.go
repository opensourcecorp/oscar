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
	cfgFileContents, readErr := Files.ReadFile(filepath.Base(t.ConfigFilePath))
	if readErr != nil {
		return fmt.Errorf("reading embedded file contents: %w", readErr)
	}

	if err := os.WriteFile(t.ConfigFilePath, cfgFileContents, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
