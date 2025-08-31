package markdownci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

type (
	markdownlintTask struct{}
)

var tasks = []ciutil.Tasker{
	markdownlintTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasMarkdown {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t markdownlintTask) InfoText() string { return "Lint (markdownlint)" }

// Run implements [ciutil.Tasker.Run].
func (t markdownlintTask) Run() error {
	if err := ciutil.RunCommand([]string{markdownlint.Name, "**/*.md"}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t markdownlintTask) Post() error { return nil }
