package markdownci

import (
	"fmt"

	nodejsci "github.com/opensourcecorp/oscar/internal/ci/nodejs"
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

// Init implements [ciutil.Tasker.Init].
func (t markdownlintTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Markdown: Installing Nodejs for markdownlint... ")

	if err := (nodejsci.BaseInitTask{}).Init(); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t markdownlintTask) Run() error {
	if err := ciutil.RunVersionedCommand(markdownlint, nil); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t markdownlintTask) Post() error { return nil }
