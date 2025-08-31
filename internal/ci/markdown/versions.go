package markdownci

import (
	"os"
	"path/filepath"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	markdownlint = ciutil.Tool{
		Name:           "markdownlint-cli2",
		ConfigFilePath: filepath.Join(os.TempDir(), ".markdownlint-cli2.yaml"),
	}
)
