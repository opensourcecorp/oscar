package shtools

import (
	"context"
	"fmt"

	"github.com/opensourcecorp/oscar/internal/tools"
)

type (
	shellcheckTask struct{}
	shfmtTask      struct{}
)

var tasks = []tools.Tasker{
	shellcheckTask{},
	shfmtTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo tools.Repo) []tools.Tasker {
	if repo.HasShell {
		return tasks
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t shellcheckTask) InfoText() string { return "Lint (shellcheck)" }

// Run implements [tools.Tasker.Run].
func (t shellcheckTask) Run(ctx context.Context) error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh`,
		shellcheck.Name,
	)}
	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t shellcheckTask) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t shfmtTask) InfoText() string { return "Format (shfmt)" }

// Run implements [tools.Tasker.Run].
func (t shfmtTask) Run(ctx context.Context) error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh`,
		shfmt.Name,
	)}
	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t shfmtTask) Post(_ context.Context) error { return nil }
