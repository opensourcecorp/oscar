package markdownci

import (
	"fmt"
	"os"
	"path/filepath"

	ciconfig "github.com/opensourcecorp/oscar/internal/ci/configfiles"
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
	cfgFileContents, err := ciconfig.Files.ReadFile(filepath.Base(markdownlint.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(markdownlint.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	args := []string{
		markdownlint.Name,
		"--config", markdownlint.ConfigFilePath,
		"**/*.md",
	}

	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t markdownlintTask) Post() error { return nil }
