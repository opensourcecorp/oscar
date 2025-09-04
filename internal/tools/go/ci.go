package igo

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/opensourcecorp/oscar/internal/tools"
	"github.com/opensourcecorp/oscar/internal/tools/toolcfg"
)

type (
	goModCheckCI   struct{}
	goFormatCI     struct{}
	generateCodeCI struct{}
	goBuildCI      struct{}
	goVetCI        struct{}
	staticcheckCI  struct{}
	reviveCI       struct{}
	errcheckCI     struct{}
	goImportsCI    struct{}
	govulncheckCI  struct{}
	goTestCI       struct{}
)

var ciTasks = []tools.Tasker{
	goModCheckCI{},
	goFormatCI{},
	generateCodeCI{},
	goBuildCI{},
	goVetCI{},
	staticcheckCI{},
	reviveCI{},
	errcheckCI{},
	goImportsCI{},
	govulncheckCI{},
	goTestCI{},
}

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

// TasksForCI returns the list of CI tasks.
func TasksForCI(repo tools.Repo) []tools.Tasker {
	if repo.HasGo {
		return ciTasks
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t goModCheckCI) InfoText() string { return "go.mod tidy check" }

// Run implements [tools.Tasker.Run].
func (t goModCheckCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "mod", "tidy"}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goModCheckCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goFormatCI) InfoText() string { return "Format" }

// Run implements [tools.Tasker.Run].
func (t goFormatCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "fmt", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goFormatCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t generateCodeCI) InfoText() string { return "Generate code" }

// Run implements [tools.Tasker.Run].
func (t generateCodeCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "generate", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t generateCodeCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goBuildCI) InfoText() string { return "Build" }

// Run implements [tools.Tasker.Run].
func (t goBuildCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "build", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goBuildCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goVetCI) InfoText() string { return "Vet" }

// Run implements [tools.Tasker.Run].
func (t goVetCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "vet", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goVetCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t staticcheckCI) InfoText() string { return "Lint (staticcheck)" }

// Run implements [tools.Tasker.Run].
func (t staticcheckCI) Run() (err error) {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(staticcheck.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(staticcheck.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	if err := goRun(staticcheck, "./..."); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t staticcheckCI) Post() error {
	if err := os.RemoveAll(staticcheck.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t reviveCI) InfoText() string { return "Lint (revive)" }

// Run implements [tools.Tasker.Run].
func (t reviveCI) Run() error {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(revive.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(revive.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	args := []string{
		"--config", revive.ConfigFilePath,
		"--set_exit_status",
		"./...",
	}

	if err := goRun(revive, args...); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t reviveCI) Post() error {
	if err := os.RemoveAll(revive.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t errcheckCI) InfoText() string { return "Lint (errcheck)" }

// Run implements [tools.Tasker.Run].
func (t errcheckCI) Run() error {
	if err := goRun(errcheck, "./..."); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t errcheckCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goImportsCI) InfoText() string { return "Format imports" }

// Run implements [tools.Tasker.Run].
func (t goImportsCI) Run() error {
	args := []string{"-l", "-w", "."}
	if err := goRun(goimports, args...); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goImportsCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t govulncheckCI) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Run implements [tools.Tasker.Run].
func (t govulncheckCI) Run() error {
	if err := goRun(govulncheck, "./..."); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t govulncheckCI) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goTestCI) InfoText() string { return "Test" }

// Run implements [tools.Tasker.Run].
func (t goTestCI) Run() error {
	if _, err := tools.RunCommand([]string{"go", "test", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goTestCI) Post() error { return nil }

// goRun is a wrapper for "go run"
func goRun(t tools.Tool, trailingArgs ...string) error {
	args := slices.Concat(
		[]string{"go", "run", fmt.Sprintf("%s@%s", t.RemotePath, t.Version)},
		trailingArgs,
	)
	if _, err := tools.RunCommand(args); err != nil {
		return fmt.Errorf("running 'go run': %w", err)
	}

	return nil
}
