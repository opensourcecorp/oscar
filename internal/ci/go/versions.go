package goci

import (
	"fmt"
	"os"
	"path/filepath"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	// NOTE: can't name this `go`, for obvious reasons
	goAsTask = ciutil.VersionedTool{
		Name: "go",
		// Placeholders are for:
		// - version
		// - kernel
		// - arch
		RemotePath:      "https://go.dev/dl/go%s.%s-%s.tar.gz",
		Version:         "1.24.4",
		VersionCheckArg: "version",
	}
	staticcheck = ciutil.VersionedTool{
		Name:           "staticcheck",
		RemotePath:     "honnef.co/go/tools/cmd/staticcheck",
		Version:        "2025.1.1",
		ConfigFilePath: filepath.Join("./staticcheck.conf"),
	}
	revive = ciutil.VersionedTool{
		Name:           "revive",
		RemotePath:     "github.com/mgechev/revive",
		Version:        "v1.11.0",
		ConfigFilePath: filepath.Join(os.TempDir(), "revive.toml"),
	}
	errcheck = ciutil.VersionedTool{
		Name:       "errcheck",
		RemotePath: "github.com/kisielk/errcheck",
		Version:    "v1.9.0",
	}
	goimports = ciutil.VersionedTool{
		Name:       "goimports",
		RemotePath: "golang.org/x/tools/cmd/goimports",
		Version:    "v0.35.0",
	}
	govulncheck = ciutil.VersionedTool{
		Name:       "govulncheck",
		RemotePath: "golang.org/x/vuln/cmd/govulncheck",
		Version:    "v1.1.4",
	}
)

// goInstall is a wrapper to run "go install" against a Go tool used for a given Task. It checks
// that the tool is up-to-date before installing.
func goInstall(vt ciutil.VersionedTool) error {
	if ciutil.IsToolUpToDate(vt) {
		return nil
	}

	args := []string{"go", "install", fmt.Sprintf("%s@%s", vt.RemotePath, vt.Version)}
	if err := ciutil.RunCommand(args); err != nil {
		return fmt.Errorf("running go install: %w", err)
	}

	return nil
}
