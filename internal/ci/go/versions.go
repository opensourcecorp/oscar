package goci

import (
	"os"
	"path/filepath"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	staticcheck = ciutil.Tool{
		Name:           "staticcheck",
		RemotePath:     "honnef.co/go/tools/cmd/staticcheck",
		Version:        "2025.1.1",
		ConfigFilePath: filepath.Join("./staticcheck.conf"),
	}
	revive = ciutil.Tool{
		Name:           "revive",
		RemotePath:     "github.com/mgechev/revive",
		Version:        "v1.11.0",
		ConfigFilePath: filepath.Join(os.TempDir(), "revive.toml"),
	}
	errcheck = ciutil.Tool{
		Name:       "errcheck",
		RemotePath: "github.com/kisielk/errcheck",
		Version:    "v1.9.0",
	}
	goimports = ciutil.Tool{
		Name:       "goimports",
		RemotePath: "golang.org/x/tools/cmd/goimports",
		Version:    "v0.35.0",
	}
	govulncheck = ciutil.Tool{
		Name:       "govulncheck",
		RemotePath: "golang.org/x/vuln/cmd/govulncheck",
		Version:    "v1.1.4",
	}
)
