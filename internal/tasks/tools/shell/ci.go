package shtools

import (
	"context"

	"github.com/opensourcecorp/oscar/internal/system"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	shellcheck struct{ taskutil.Tool }
	shfmt      struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasShell {
		return []taskutil.Tasker{
			shellcheck{
				Tool: taskutil.Tool{
					RunArgs: []string{"bash", "-c", `
						shopt -s globstar
						ls **/*.sh || exit 0
						shellcheck **/*.sh`,
					},
				},
			},
			shfmt{
				Tool: taskutil.Tool{
					RunArgs: []string{"bash", "-c", `
						shopt -s globstar
						ls **/*.sh || exit 0
						shfmt **/*.sh`,
					},
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t shellcheck) InfoText() string { return "Lint (shellcheck)" }

// Run implements [taskutil.Tasker.Run].
func (t shellcheck) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t shellcheck) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t shfmt) InfoText() string { return "Format (shfmt)" }

// Run implements [taskutil.Tasker.Run].
func (t shfmt) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t shfmt) Post(_ context.Context) error { return nil }
