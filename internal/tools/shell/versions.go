package shtools

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
)
