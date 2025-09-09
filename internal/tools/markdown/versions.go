package mdtools

import (
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tools"
)

var (
	markdownlint = tools.Tool{
		Name:           "markdownlint-cli2",
		ConfigFilePath: filepath.Join(os.TempDir(), ".markdownlint-cli2.yaml"),
	}
)
