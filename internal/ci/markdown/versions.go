package markdownci

import ciutil "github.com/opensourcecorp/oscar/internal/ci/util"

var (
	markdownlint = ciutil.VersionedTool{
		Name:       "markdownlint-cli2",
		Version:    "v0.18.1",
		RunCommand: []string{"npx"},
	}
)
