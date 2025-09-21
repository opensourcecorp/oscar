package pytools

import (
	"context"

	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	buildTask  struct{ taskutil.Tool }
	ruffLint   struct{ taskutil.Tool }
	ruffFormat struct{ taskutil.Tool }
	pydoclint  struct{ taskutil.Tool }
	mypy       struct{ taskutil.Tool }
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
					RunArgs: []string{"ruff", "check", "--fix", "./src"},
				},
			},
			ruffFormat{
				Tool: taskutil.Tool{
					RunArgs: []string{"ruff", "format", "./src"},
				},
			},
			pydoclint{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "pydoclint", "./src"},
				},
			},
			mypy{
				Tool: taskutil.Tool{
					RunArgs: []string{"uvx", "mypy", "./src"},
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t buildTask) InfoText() string { return "Build" }

// Run implements [taskutil.Tasker.Run].
func (t buildTask) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t buildTask) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t ruffLint) InfoText() string { return "Lint (ruff)" }

// Run implements [taskutil.Tasker.Run].
func (t ruffLint) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ruffLint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t ruffFormat) InfoText() string { return "Format (ruff)" }

// Run implements [taskutil.Tasker.Run].
func (t ruffFormat) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ruffFormat) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t pydoclint) InfoText() string { return "Lint (pydoclint)" }

// Run implements [taskutil.Tasker.Run].
func (t pydoclint) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t pydoclint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t mypy) InfoText() string { return "Type-check (mypy)" }

// Run implements [taskutil.Tasker.Run].
func (t mypy) Exec(ctx context.Context) error {
	if _, err := taskutil.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t mypy) Post(_ context.Context) error { return nil }
