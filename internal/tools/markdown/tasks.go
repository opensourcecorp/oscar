package mdtools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tools"
	"github.com/opensourcecorp/oscar/internal/tools/toolcfg"
)

type (
	markdownlintTask struct{}
)

var tasks = []tools.Tasker{
	markdownlintTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo tools.Repo) []tools.Tasker {
	if repo.HasMarkdown {
		return tasks
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t markdownlintTask) InfoText() string { return "Lint (markdownlint)" }

// Run implements [tools.Tasker.Run].
func (t markdownlintTask) Run(ctx context.Context) error {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(markdownlint.ConfigFilePath))
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

	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t markdownlintTask) Post(_ context.Context) error { return nil }
