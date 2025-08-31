package igo

import (
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tools"
)

var (
	staticcheck = tools.Tool{
		Name:           "staticcheck",
		RemotePath:     "honnef.co/go/tools/cmd/staticcheck",
		Version:        "2025.1.1",
		ConfigFilePath: filepath.Join("./staticcheck.conf"),
	}
	revive = tools.Tool{
		Name:           "revive",
		RemotePath:     "github.com/mgechev/revive",
		Version:        "v1.11.0",
		ConfigFilePath: filepath.Join(os.TempDir(), "revive.toml"),
	}
	errcheck = tools.Tool{
		Name:       "errcheck",
		RemotePath: "github.com/kisielk/errcheck",
		Version:    "v1.9.0",
	}
	goimports = tools.Tool{
		Name:       "goimports",
		RemotePath: "golang.org/x/tools/cmd/goimports",
		Version:    "v0.35.0",
	}
	govulncheck = tools.Tool{
		Name:       "govulncheck",
		RemotePath: "golang.org/x/vuln/cmd/govulncheck",
		Version:    "v1.1.4",
	}
)
