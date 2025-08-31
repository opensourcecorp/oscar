package shellci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	shellcheck = ciutil.Tool{
		Name:    "shellcheck",
		Version: "v0.11.0",
	}
	shfmt = ciutil.Tool{
		Name:    "shfmt",
		Version: "v3.12.0",
	}
	bats = ciutil.Tool{
		Name: "bats",
		// NOTE: bats just gets cloned then installed with its install script
		RemotePath: "https://github.com/bats-core/bats-core.git",
		Version:    "v1.12.0",
	}
)
