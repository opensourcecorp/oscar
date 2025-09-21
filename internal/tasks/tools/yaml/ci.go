package yamltools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tasks/tools/toolcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	yamlfmt  struct{ taskutil.Tool }
	yamllint struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasYaml {
		return []taskutil.Tasker{
			yamlfmt{
				Tool: taskutil.Tool{
					RunArgs: []string{"bash", "-c",
						fmt.Sprintf(
							`yamlfmt -conf {{ConfigFilePath}} $(%s)`,
							taskutil.GetFileTypeListerCommand("yaml"),
						),
					},
					ConfigFilePath: filepath.Join(os.TempDir(), ".yamlfmt"),
				},
			},
			yamllint{
				Tool: taskutil.Tool{
					RunArgs: []string{"bash", "-c",
						fmt.Sprintf(
							`yamllint --strict --config-file {{ConfigFilePath}} $(%s)`,
							taskutil.GetFileTypeListerCommand("yaml"),
						),
					},
					ConfigFilePath: filepath.Join(os.TempDir(), ".yamllint"),
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t yamllint) InfoText() string { return "Lint (yamllint)" }

// Run implements [taskutil.Tasker.Run].
func (t yamllint) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t yamllint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t yamlfmt) InfoText() string { return "Format (yamlfmt)" }

// Run implements [taskutil.Tasker.Run].
func (t yamlfmt) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t yamlfmt) Post(_ context.Context) error { return nil }
