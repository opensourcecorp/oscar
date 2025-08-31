package goci

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	ciconfig "github.com/opensourcecorp/oscar/internal/ci/configfiles"
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

// A list of tasks that all implement [ciutil.Tasker], for Go.
type (
	goModCheckTask   struct{}
	goFormatTask     struct{}
	generateCodeTask struct{}
	goBuildTask      struct{}
	goVetTask        struct{}
	staticcheckTask  struct{}
	reviveTask       struct{}
	errcheckTask     struct{}
	goImportsTask    struct{}
	govulncheckTask  struct{}
	goTestTask       struct{}
)

var tasks = []ciutil.Tasker{
	goModCheckTask{},
	goFormatTask{},
	generateCodeTask{},
	goBuildTask{},
	goVetTask{},
	staticcheckTask{},
	reviveTask{},
	errcheckTask{},
	goImportsTask{},
	govulncheckTask{},
	goTestTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasGo {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t goModCheckTask) InfoText() string { return "go.mod tidy check" }

// Run implements [ciutil.Tasker.Run].
func (t goModCheckTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "mod", "tidy"}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goModCheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goFormatTask) InfoText() string { return "Format" }

// Run implements [ciutil.Tasker.Run].
func (t goFormatTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "fmt", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goFormatTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t generateCodeTask) InfoText() string { return "Generate code" }

// Run implements [ciutil.Tasker.Run].
func (t generateCodeTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "generate", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t generateCodeTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goBuildTask) InfoText() string { return "Build" }

// Run implements [ciutil.Tasker.Run].
func (t goBuildTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "build", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goBuildTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goVetTask) InfoText() string { return "Vet" }

// Run implements [ciutil.Tasker.Run].
func (t goVetTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "vet", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goVetTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t staticcheckTask) InfoText() string { return "Lint (staticcheck)" }

// Run implements [ciutil.Tasker.Run].
func (t staticcheckTask) Run() (err error) {
	cfgFileContents, err := ciconfig.Files.ReadFile(filepath.Base(staticcheck.ConfigFilePath))
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

// Post implements [ciutil.Tasker.Post].
func (t staticcheckTask) Post() error {
	if err := os.RemoveAll(staticcheck.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t reviveTask) InfoText() string { return "Lint (revive)" }

// Run implements [ciutil.Tasker.Run].
func (t reviveTask) Run() error {
	cfgFileContents, err := ciconfig.Files.ReadFile(filepath.Base(revive.ConfigFilePath))
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

// Post implements [ciutil.Tasker.Post].
func (t reviveTask) Post() error {
	if err := os.RemoveAll(revive.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t errcheckTask) InfoText() string { return "Lint (errcheck)" }

// Run implements [ciutil.Tasker.Run].
func (t errcheckTask) Run() error {
	if err := goRun(errcheck, "./..."); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t errcheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goImportsTask) InfoText() string { return "Format imports" }

// Run implements [ciutil.Tasker.Run].
func (t goImportsTask) Run() error {
	args := []string{"-l", "-w", "."}
	if err := goRun(goimports, args...); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goImportsTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t govulncheckTask) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Run implements [ciutil.Tasker.Run].
func (t govulncheckTask) Run() error {
	if err := goRun(govulncheck, "./..."); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t govulncheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goTestTask) InfoText() string { return "Test" }

// Run implements [ciutil.Tasker.Run].
func (t goTestTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "test", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goTestTask) Post() error { return nil }

// goRun is a wrapper for "go run"
func goRun(t ciutil.Tool, trailingArgs ...string) error {
	args := slices.Concat(
		[]string{"go", "run", fmt.Sprintf("%s@%s", t.RemotePath, t.Version)},
		trailingArgs,
	)
	if err := ciutil.RunCommand(args); err != nil {
		return fmt.Errorf("running 'go run': %w", err)
	}

	return nil
}
