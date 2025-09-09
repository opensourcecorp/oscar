package yamltools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tools"
	"github.com/opensourcecorp/oscar/internal/tools/toolcfg"
)

type (
	yamlfmtTask  struct{}
	yamllintTask struct{}
)

var tasks = []tools.Tasker{
	yamlfmtTask{},
	yamllintTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo tools.Repo) []tools.Tasker {
	if repo.HasYaml {
		return tasks
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t yamllintTask) InfoText() string { return "Lint (yamllint)" }

// Run implements [tools.Tasker.Run].
func (t yamllintTask) Run(ctx context.Context) error {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(yamllint.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(yamllint.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	args := []string{"bash", "-c",
		fmt.Sprintf(
			`%s --strict --config-file %s $(%s)`,
			yamllint.Name, yamllint.ConfigFilePath, tools.GetFileTypeListerCommand("yaml"),
		),
	}

	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t yamllintTask) Post(_ context.Context) error { return nil }

// InfoText implements [tools.Tasker.InfoText].
func (t yamlfmtTask) InfoText() string { return "Format (yamlfmt)" }

// Run implements [tools.Tasker.Run].
func (t yamlfmtTask) Run(ctx context.Context) error {
	cfgFileContents, err := toolcfg.Files.ReadFile(filepath.Base(yamlfmt.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(yamlfmt.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	args := []string{"bash", "-c",
		fmt.Sprintf(
			`%s -conf %s $(%s)`,
			yamlfmt.Name, yamlfmt.ConfigFilePath, tools.GetFileTypeListerCommand("yaml"),
		),
	}

	if _, err := tools.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t yamlfmtTask) Post(_ context.Context) error { return nil }
