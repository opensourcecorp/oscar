package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tasks/tools/toolcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	goModCheck     struct{ taskutil.Tool }
	goFormat       struct{ taskutil.Tool }
	generateCodeCI struct{ taskutil.Tool }
	goBuildCI      struct{ taskutil.Tool }
	goVet          struct{ taskutil.Tool }
	staticcheck    struct{ taskutil.Tool }
	revive         struct{ taskutil.Tool }
	errcheck       struct{ taskutil.Tool }
	goImports      struct{ taskutil.Tool }
	govulncheck    struct{ taskutil.Tool }
	goTest         struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasGo {
		return []taskutil.Tasker{
			goModCheck{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "mod", "tidy"},
				},
			},
			goFormat{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "fmt", "./..."},
				},
			},
			goImports{
				Tool: taskutil.Tool{
					RunArgs: []string{"goimports", "-l", "-w", "."},
				},
			},
			generateCodeCI{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "generate", "./..."},
				},
			},
			goBuildCI{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "build", "./..."},
				},
			},
			goVet{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "vet", "./..."},
				},
			},
			staticcheck{
				Tool: taskutil.Tool{
					RunArgs: []string{"staticcheck", "./..."},
					// NOTE: staticcheck does not have a flag to point to a config file, so we need
					// to put it at the repo root
					ConfigFilePath: filepath.Join("staticcheck.conf"),
				},
			},
			revive{
				Tool: taskutil.Tool{
					RunArgs: []string{
						"revive", "--config", "{{ConfigFilePath}}", "--set_exit_status", "./...",
					},
					ConfigFilePath: filepath.Join(os.TempDir(), "revive.toml"),
				},
			},
			errcheck{
				Tool: taskutil.Tool{
					RunArgs: []string{"errcheck", "./..."},
				},
			},
			govulncheck{
				Tool: taskutil.Tool{
					RunArgs: []string{"govulncheck", "./..."},
				},
			},
			goTest{
				Tool: taskutil.Tool{
					RunArgs: []string{"go", "test", "./..."},
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t goModCheck) InfoText() string { return "go.mod tidy check" }

// Exec implements [taskutil.Tasker.Exec].
func (t goModCheck) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goModCheck) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t goFormat) InfoText() string { return "Format" }

// Exec implements [taskutil.Tasker.Exec].
func (t goFormat) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goFormat) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t goImports) InfoText() string { return "Format imports" }

// Exec implements [taskutil.Tasker.Exec].
func (t goImports) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goImports) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t generateCodeCI) InfoText() string { return "Generate code" }

// Exec implements [taskutil.Tasker.Exec].
func (t generateCodeCI) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	// Generating code will likely throw diffs if not also addressing other formatting CI checks, so
	// run those here again as well.
	tasks := NewTasksForCI(taskutil.Repo{HasGo: true})
	for _, task := range tasks {
		if wantToRun, isOfType := task.(goFormat); isOfType {
			if err := wantToRun.Exec(ctx); err != nil {
				return fmt.Errorf("running Go formatter after code generation: %w", err)
			}
		}
		if wantToRun, isOfType := task.(goImports); isOfType {
			if err := wantToRun.Exec(ctx); err != nil {
				return fmt.Errorf("running Go formatter after code generation: %w", err)
			}
		}
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t generateCodeCI) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t goBuildCI) InfoText() string { return "Build" }

// Exec implements [taskutil.Tasker.Exec].
func (t goBuildCI) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goBuildCI) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t goVet) InfoText() string { return "Vet" }

// Exec implements [taskutil.Tasker.Exec].
func (t goVet) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goVet) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t staticcheck) InfoText() string { return "Lint (staticcheck)" }

// Exec implements [taskutil.Tasker.Exec].
func (t staticcheck) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t staticcheck) Post(_ context.Context) error {
	if err := os.RemoveAll(t.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t revive) InfoText() string { return "Lint (revive)" }

// Exec implements [taskutil.Tasker.Exec].
func (t revive) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t revive) Post(_ context.Context) error {
	if err := os.RemoveAll(t.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t errcheck) InfoText() string { return "Lint (errcheck)" }

// Exec implements [taskutil.Tasker.Exec].
func (t errcheck) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t errcheck) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t govulncheck) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Exec implements [taskutil.Tasker.Exec].
func (t govulncheck) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t govulncheck) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t goTest) InfoText() string { return "Tests" }

// Exec implements [taskutil.Tasker.Exec].
func (t goTest) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goTest) Post(_ context.Context) error { return nil }
