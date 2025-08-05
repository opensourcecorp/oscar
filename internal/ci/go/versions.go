package goci

import (
	"fmt"
	"os/exec"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	staticcheck = ciutil.VersionedTask{
		Name:       "staticcheck",
		RemotePath: "honnef.co/go/tools/cmd/staticcheck",
		Version:    "2025.1.1",
	}
	revive = ciutil.VersionedTask{
		Name:       "revive",
		RemotePath: "github.com/mgechev/revive",
		Version:    "v1.11.0",
	}
	errcheck = ciutil.VersionedTask{
		Name:       "errcheck",
		RemotePath: "github.com/kisielk/errcheck",
		Version:    "v1.9.0",
	}
	goimports = ciutil.VersionedTask{
		Name:       "goimports",
		RemotePath: "golang.org/x/tools/cmd/goimports",
		Version:    "v0.35.0",
	}
	govulncheck = ciutil.VersionedTask{
		Name:       "govulncheck",
		RemotePath: "golang.org/x/vuln/cmd/govulncheck",
		Version:    "v1.1.4",
	}
)

// TODO:
func goInstall(i ciutil.VersionedTask) error {
	cmd := exec.Command("go", "install", fmt.Sprintf("%s@%s", i.RemotePath, i.Version))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(
			"running go install for '%s': %w -- output:\n%s",
			i.Name, err, string(output),
		)
	}

	return nil
}
