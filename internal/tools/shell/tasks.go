package shell

import (
	"fmt"

	"github.com/opensourcecorp/oscar/internal/tools"
)

type (
	shellcheckTask struct{}
	shfmtTask      struct{}
	batsTask       struct{}
)

var tasks = []tools.Tasker{
	shellcheckTask{},
	shfmtTask{},
	batsTask{},
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
func (t shellcheckTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shellcheck.Name)}
	if err := tools.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t shellcheckTask) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t shfmtTask) InfoText() string { return "Format (shfmt)" }

// Run implements [tools.Tasker.Run].
func (t shfmtTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shfmt.Name)}
	if err := tools.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t shfmtTask) Post() error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t batsTask) InfoText() string { return "Test (bats)" }

// Run implements [tools.Tasker.Run].
func (t batsTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		# Don't run if no bats files found, otherwise it will error out
		ls **/*.bats || exit 0
		%s **/*.bats
	`, bats.Name)}
	if err := tools.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t batsTask) Post() error { return nil }
