package markdownci

import (
	"os"
	"path/filepath"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	markdownlint = ciutil.Tool{
		Name: "markdownlint-cli2",
		// NOTE: the tool requires that, regardless of location, that you still give it a name
		// pattern it expects, hence the verbose one used here
		ConfigFilePath: filepath.Join(os.TempDir(), ".markdownlint-cli2.yaml"),
	}
)
