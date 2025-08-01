package ci

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

func initGo() error {
	// Adding the default GOPATH/bin is easier than trying to figure out a user's own custom
	// settings
	_ = os.Setenv(
		"PATH",
		fmt.Sprintf("%s:%s", filepath.Join(os.Getenv("HOME"), "go", "bin"), os.Getenv("PATH")),
	)

	iprint.Debugf("PATH after Go init: %s\n", os.Getenv("PATH"))

	installCommands := [][]string{
		{"go", "install",
			fmt.Sprintf("honnef.co/go/tools/cmd/staticcheck@%s", consts.GoCIVersions.Staticcheck)},
		{"go", "install",
			fmt.Sprintf("github.com/mgechev/revive@%s", consts.GoCIVersions.Revive)},
		{"go", "install",
			fmt.Sprintf("github.com/kisielk/errcheck@%s", consts.GoCIVersions.Errcheck)},
		{"go", "install",
			fmt.Sprintf("golang.org/x/tools/cmd/goimports@%s", consts.GoCIVersions.Goimports)},
		{"go", "install",
			fmt.Sprintf("golang.org/x/vuln/cmd/govulncheck@%s", consts.GoCIVersions.Govulncheck)},
	}
	for _, i := range installCommands {
		if err := runInitCommand(i); err != nil {
			return err
		}
	}

	return nil
}

func getGoConfigs(repo Repo) []Config {
	if repo.HasGo {
		return []Config{
			{
				LanguageName: "Go",
				Tasks: []Task{
					{
						InfoText: "Init",
						InitFunc: initGo,
					},
					{
						InfoText:  "Go mod check",
						RunScript: []string{"go", "mod", "tidy"},
					},
					{
						InfoText:  "Format code",
						RunScript: []string{"go", "fmt", "./..."},
					},
					{
						InfoText:  "Generate code",
						RunScript: []string{"bash", "-c", "go generate ./... && go fmt ./..."},
					},
					{
						InfoText:  "Build",
						RunScript: []string{"go", "build", "./..."},
					},
					{
						InfoText:  "Vet",
						RunScript: []string{"go", "vet", "./..."},
					},
					{
						InfoText:  "Linter (staticcheck)",
						RunScript: []string{"staticcheck", "./..."},
					},
					{
						InfoText:  "Linter (revive)",
						RunScript: []string{"revive", "--set_exit_status", "./..."},
					},
					{
						InfoText:  "Linter (errcheck)",
						RunScript: []string{"go", "run", "github.com/kisielk/errcheck@latest", "./..."},
					},
					{
						InfoText:  "Linter (goimports)",
						RunScript: []string{"go", "run", "golang.org/x/tools/cmd/goimports@latest", "-l", "-w", "."},
					},
					{
						InfoText:  "Vulnerability scanner (govulncheck)",
						RunScript: []string{"go", "run", "golang.org/x/vuln/cmd/govulncheck@latest", "./..."},
					},
					{
						InfoText:  "Run tests",
						RunScript: []string{"go", "test", "./..."},
					},
				},
			},
		}
	}

	return nil
}
