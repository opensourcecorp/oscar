package pytools

import (
	"context"
	"fmt"

	"github.com/opensourcecorp/oscar/internal/system"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	buildTask  struct{ taskutil.Tool }
	ruffLint   struct{ taskutil.Tool }
	ruffFormat struct{ taskutil.Tool }
	pydoclint  struct{ taskutil.Tool }
	ty         struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasPython {
		return []taskutil.Tasker{
			buildTask{
				Tool: taskutil.Tool{
					RunArgs: []string{"uv", "build"},
				},
			},
			ruffLint{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "ruff", "check", "--fix", "./src"},
				},
			},
			ruffFormat{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "ruff", "format", "./src"},
				},
			},
			pydoclint{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "pydoclint", "./src"},
				},
			},
			ty{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "ty", "check"},
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t buildTask) InfoText() string { return "Build" }

// Exec implements [taskutil.Tasker.Exec].
func (t buildTask) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return fmt.Errorf("running build: %w", err)
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t buildTask) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t ruffLint) InfoText() string { return "Lint (ruff)" }

// Exec implements [taskutil.Tasker.Exec].
func (t ruffLint) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return fmt.Errorf("running ruff linter: %w", err)
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ruffLint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t ruffFormat) InfoText() string { return "Format (ruff)" }

// Exec implements [taskutil.Tasker.Exec].
func (t ruffFormat) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return fmt.Errorf("running ruff formatter: %w", err)
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ruffFormat) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t pydoclint) InfoText() string { return "Lint (pydoclint)" }

// Exec implements [taskutil.Tasker.Exec].
func (t pydoclint) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return fmt.Errorf("running pydoclint: %w", err)
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t pydoclint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t ty) InfoText() string { return "Type-check (ty)" }

// Exec implements [taskutil.Tasker.Exec].
func (t ty) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return fmt.Errorf("running ty type checker: %w", err)
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ty) Post(_ context.Context) error { return nil }
