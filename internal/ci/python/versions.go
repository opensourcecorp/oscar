package pythonci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	// NOTE: even though ruff is used for both linting & formatting, their implementations differ,
	// so we need two distinct [ciutil.Tool]s.
	ruffLint = ciutil.Tool{
		Name:    "ruff",
		Version: "0.12.7",
	}
	ruffFormat = ciutil.Tool{
		Name:    "ruff",
		Version: "0.12.7",
	}
	pydoclint = ciutil.Tool{
		Name:    "pydoclint",
		Version: "0.6.6",
	}
	mypy = ciutil.Tool{
		Name:    "mypy",
		Version: "1.17.1",
	}
)
