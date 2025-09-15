package mdtools

import (
	"context"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tasks/tools/toolcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	markdownlint struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasMarkdown {
		return []taskutil.Tasker{
			markdownlint{
				Tool: taskutil.Tool{
					RunArgs:        []string{"markdownlint-cli2", "--config", "{{ConfigFilePath}}", "**/*.md"},
					ConfigFilePath: filepath.Join(os.TempDir(), ".markdownlint-cli2.yaml"),
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t markdownlint) InfoText() string { return "Lint (markdownlint)" }

// Exec implements [taskutil.Tasker.Exec].
func (t markdownlint) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t markdownlint) Post(_ context.Context) error { return nil }
