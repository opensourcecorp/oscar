package pytools

import (
	"github.com/opensourcecorp/oscar/internal/tools"
)

var (
	// NOTE: even though ruff is used for both linting & formatting, their implementations differ,
	// so we need two distinct [ciutil.Tool]s.
	ruffLint = tools.Tool{
		Name:    "ruff",
		Version: "0.12.7",
	}
	ruffFormat = tools.Tool{
		Name:    "ruff",
		Version: "0.12.7",
	}
	pydoclint = tools.Tool{
		Name:    "pydoclint",
		Version: "0.6.6",
	}
	mypy = tools.Tool{
		Name:    "mypy",
		Version: "1.17.1",
	}
)
