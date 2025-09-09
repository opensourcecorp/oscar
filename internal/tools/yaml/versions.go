package yamltools

import (
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tools"
)

var (
	yamlfmt = tools.Tool{
		Name:           "yamlfmt",
		ConfigFilePath: filepath.Join(os.TempDir(), ".yamlfmt"),
	}
	yamllint = tools.Tool{
		Name:           "yamllint",
		ConfigFilePath: filepath.Join(os.TempDir(), ".yamllint"),
	}
)
