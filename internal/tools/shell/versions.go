package shell

import (
	"github.com/opensourcecorp/oscar/internal/tools"
)

var (
	shellcheck = tools.Tool{
		Name:    "shellcheck",
		Version: "v0.11.0",
	}
	shfmt = tools.Tool{
		Name:    "shfmt",
		Version: "v3.12.0",
	}
	bats = tools.Tool{
		Name: "bats",
		// NOTE: bats just gets cloned then installed with its install script
		RemotePath: "https://github.com/bats-core/bats-core.git",
		Version:    "v1.12.0",
	}
)
