package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/system"
	"github.com/opensourcecorp/oscar/internal/tasks/tools/toolcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	goModCheck     struct{ taskutil.Tool }
	generateCodeCI struct{ taskutil.Tool }
	goBuildCI      struct{ taskutil.Tool }
	golangciFmt    struct{ taskutil.Tool }
	golangciLint   struct{ taskutil.Tool }
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
			golangciFmt{
				Tool: taskutil.Tool{
					RunArgs: []string{
						"golangci-lint", "fmt", "--config", "{{ConfigFilePath}}",
					},
					ConfigFilePath: filepath.Join(os.TempDir(), ".golangci.yaml"),
				},
			},
			golangciLint{
				Tool: taskutil.Tool{
					RunArgs: []string{
						"golangci-lint", "run", "--config", "{{ConfigFilePath}}",
					},
					ConfigFilePath: filepath.Join(os.TempDir(), ".golangci.yaml"),
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
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goModCheck) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t generateCodeCI) InfoText() string { return "Generate code" }

// Exec implements [taskutil.Tasker.Exec].
func (t generateCodeCI) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	// Generating code will likely throw diffs if not also addressing other formatting CI checks, so
	// run those here as well.
	tasks := NewTasksForCI(taskutil.Repo{HasGo: true})
	for _, task := range tasks {
		if wantToRun, isOfType := task.(golangciFmt); isOfType {
			err := wantToRun.Exec(ctx)
			if err != nil {
				return fmt.Errorf("running Go formatters after code generation: %w", err)
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
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goBuildCI) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t golangciFmt) InfoText() string { return "Format (golangci-lint)" }

// Exec implements [taskutil.Tasker.Exec].
func (t golangciFmt) Exec(ctx context.Context) error {
	err := toolcfg.SetupConfigFile(t.Tool)
	if err != nil {
		return err
	}

	if _, err := system.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t golangciFmt) Post(_ context.Context) error {
	err := os.RemoveAll(t.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t golangciLint) InfoText() string { return "Lint (golangci-lint)" }

// Exec implements [taskutil.Tasker.Exec].
func (t golangciLint) Exec(ctx context.Context) error {
	err := toolcfg.SetupConfigFile(t.Tool)
	if err != nil {
		return err
	}

	if _, err := system.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t golangciLint) Post(_ context.Context) error {
	err := os.RemoveAll(t.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t govulncheck) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Exec implements [taskutil.Tasker.Exec].
func (t govulncheck) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
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
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t goTest) Post(_ context.Context) error { return nil }
