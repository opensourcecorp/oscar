package nodejsci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	// runCmd = []string{"npx"}

	nodejs = ciutil.VersionedTool{
		Name: "node",
		// Placeholders are for:
		// - version
		// - version (again)
		// - kernel
		// - arch
		// - file extension (they release macOS as .tar.gz, and Linux as .tar.xz, for some reason)
		RemotePath: "https://nodejs.org/dist/%s/node-%s-%s-%s%s",
		Version:    "v22.18.0",
	}
)
