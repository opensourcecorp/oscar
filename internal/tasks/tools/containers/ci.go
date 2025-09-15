package containertools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/tasks/tools/toolcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	hadolint struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasContainerfile {
		return []taskutil.Tasker{
			hadolint{
				Tool: taskutil.Tool{
					RunArgs: []string{"bash", "-c",
						fmt.Sprintf(
							`hadolint --config {{ConfigFilePath}} $(%s)`,
							taskutil.GetFileTypeListerCommand("containerfile"),
						),
					},
					ConfigFilePath: filepath.Join(os.TempDir(), "hadolint.yaml"),
				},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t hadolint) InfoText() string { return "Lint (hadolint)" }

// Run implements [taskutil.Tasker.Run].
func (t hadolint) Exec(ctx context.Context) error {
	if err := toolcfg.SetupConfigFile(t.Tool); err != nil {
		return err
	}

	if _, err := taskutil.RunCommand(ctx, t.RenderRunCommandArgs()); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t hadolint) Post(_ context.Context) error { return nil }
