package shellci

import (
	"fmt"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

type (
	shellcheckTask struct{}
	shfmtTask      struct{}
	batsTask       struct{}
)

var tasks = []ciutil.Tasker{
	shellcheckTask{},
	shfmtTask{},
	batsTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasShell {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t shellcheckTask) InfoText() string { return "Lint (shellcheck)" }

// Run implements [ciutil.Tasker.Run].
func (t shellcheckTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shellcheck.Name)}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t shellcheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t shfmtTask) InfoText() string { return "Format (shfmt)" }

// Run implements [ciutil.Tasker.Run].
func (t shfmtTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shfmt.Name)}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t shfmtTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t batsTask) InfoText() string { return "Test (bats)" }

// Run implements [ciutil.Tasker.Run].
func (t batsTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		# Don't run if no bats files found, otherwise it will error out
		ls **/*.bats || exit 0
		%s **/*.bats
	`, bats.Name)}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t batsTask) Post() error { return nil }
