package pythonci

import (
	"fmt"
	"slices"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

type (
	baseConfigTask struct{}
	buildTask      struct{}
	ruffLintTask   struct{}
	ruffFormatTask struct{}
	pydoclintTask  struct{}
	mypyTask       struct{}
)

var tasks = []ciutil.Tasker{
	baseConfigTask{},
	buildTask{},
	ruffLintTask{},
	ruffFormatTask{},
	pydoclintTask{},
	mypyTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasPython {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t baseConfigTask) InfoText() string { return "" }

// Run implements [ciutil.Tasker.Run].
func (t baseConfigTask) Run() error {
	// ciutil.PlaceConfigFile("pyproject.toml")

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t baseConfigTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t buildTask) InfoText() string { return "Build" }

// Run implements [ciutil.Tasker.Run].
func (t buildTask) Run() error {
	if err := ciutil.RunCommand([]string{"uv", "build"}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t buildTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t ruffLintTask) InfoText() string { return "Lint (ruff)" }

// Run implements [ciutil.Tasker.Run].
func (t ruffLintTask) Run() error {
	if err := pyRun(ruffLint, "check", "--fix", "./src"); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t ruffLintTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t ruffFormatTask) InfoText() string { return "Format (ruff)" }

// Run implements [ciutil.Tasker.Run].
func (t ruffFormatTask) Run() error {
	if err := pyRun(ruffFormat, "format", "./src"); err != nil {
		return err
	}
	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t ruffFormatTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t pydoclintTask) InfoText() string { return "Lint (pydoclint)" }

// Run implements [ciutil.Tasker.Run].
func (t pydoclintTask) Run() error {
	if err := pyRun(pydoclint, "./src"); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t pydoclintTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t mypyTask) InfoText() string { return "Type-check (mypy)" }

// Run implements [ciutil.Tasker.Run].
func (t mypyTask) Run() error {
	if err := pyRun(mypy, "./src"); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t mypyTask) Post() error { return nil }

// pyRun is a wrapper for "uvx"
func pyRun(t ciutil.Tool, trailingArgs ...string) error {
	args := slices.Concat(
		[]string{"uvx", fmt.Sprintf("%s@%s", t.Name, t.Version)},
		trailingArgs,
	)
	if err := ciutil.RunCommand(args); err != nil {
		return fmt.Errorf("running 'uvx': %w", err)
	}

	return nil
}
