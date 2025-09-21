package gittagtools

import (
	"context"
	"fmt"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	createAndPushTag struct{ taskutil.Tool }
)

// NewTasksForDelivery returns the list of Delivery tasks.
func NewTasksForDelivery(_ taskutil.Repo) ([]taskutil.Tasker, error) {
	out := []taskutil.Tasker{
		createAndPushTag{},
	}

	return out, nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t createAndPushTag) InfoText() string { return "Create & Push Git Tag" }

// Exec implements [taskutil.Tasker.Exec].
func (t createAndPushTag) Exec(ctx context.Context) error {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return err
	}

	args := []string{"bash", "-c", fmt.Sprintf(`
		git tag v%s
		git push --tags
	`, cfg.GetVersion(),
	)}

	if _, err := taskutil.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t createAndPushTag) Post(_ context.Context) error { return nil }
