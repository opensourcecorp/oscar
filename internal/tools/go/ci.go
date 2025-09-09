package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
		ConfigFilePath: filepath.Join("./staticcheck.conf"),
	}
	revive = tools.Tool{
		Name:           "revive",
		ConfigFilePath: filepath.Join(os.TempDir(), "revive.toml"),
	}
	errcheck = tools.Tool{
		Name: "errcheck",
	}
	goimports = tools.Tool{
		Name: "goimports",
	}
	govulncheck = tools.Tool{
		Name: "govulncheck",
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
func (t goModCheckCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "mod", "tidy"}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goModCheckCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goFormatCI) InfoText() string { return "Format" }

// Run implements [tools.Tasker.Run].
func (t goFormatCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "fmt", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goFormatCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t generateCodeCI) InfoText() string { return "Generate code" }

// Run implements [tools.Tasker.Run].
func (t generateCodeCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "generate", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t generateCodeCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goBuildCI) InfoText() string { return "Build" }

// Run implements [tools.Tasker.Run].
func (t goBuildCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "build", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goBuildCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goVetCI) InfoText() string { return "Vet" }

// Run implements [tools.Tasker.Run].
func (t goVetCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "vet", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goVetCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t staticcheckCI) InfoText() string { return "Lint (staticcheck)" }

// Run implements [tools.Tasker.Run].
func (t staticcheckCI) Run(ctx context.Context) (err error) {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(staticcheck.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(staticcheck.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	if _, err := tools.RunCommand(ctx, []string{staticcheck.Name, "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t staticcheckCI) Post(_ context.Context) error {
	if err := os.RemoveAll(staticcheck.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t reviveCI) InfoText() string { return "Lint (revive)" }

// Run implements [tools.Tasker.Run].
func (t reviveCI) Run(ctx context.Context) error {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(revive.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(revive.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	args := []string{
		revive.Name,
		"--config", revive.ConfigFilePath,
		"--set_exit_status",
		"./...",
	}

	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t reviveCI) Post(_ context.Context) error {
	if err := os.RemoveAll(revive.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t errcheckCI) InfoText() string { return "Lint (errcheck)" }

// Run implements [tools.Tasker.Run].
func (t errcheckCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{errcheck.Name, "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t errcheckCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goImportsCI) InfoText() string { return "Format imports" }

// Run implements [tools.Tasker.Run].
func (t goImportsCI) Run(ctx context.Context) error {
	args := []string{goimports.Name, "-l", "-w", "."}
	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goImportsCI) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t govulncheckCI) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Run implements [tools.Tasker.Run].
func (t govulncheckCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{govulncheck.Name, "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t govulncheckCI) Post(ctx context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t goTestCI) InfoText() string { return "Test" }

// Run implements [tools.Tasker.Run].
func (t goTestCI) Run(ctx context.Context) error {
	if _, err := tools.RunCommand(ctx, []string{"go", "test", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t goTestCI) Post(ctx context.Context) error { return nil }
