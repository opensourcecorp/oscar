package tftools

import (
	"context"
	"errors"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	tfApply struct{ taskutil.Tool }
)

// NewTasksForDelivery returns the list of Delivery tasks.
func NewTasksForDelivery(repo taskutil.Repo) ([]taskutil.Tasker, error) {
	if repo.HasTerraform {
		return []taskutil.Tasker{tfApply{}}, nil
	}

	return nil, nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t tfApply) InfoText() string { return "Terraform Apply" }

// Exec implements [taskutil.Tasker.Exec].
func (t tfApply) Exec(ctx context.Context) error {
	_, err := oscarcfg.Get()
	if err != nil {
		return err
	}

	return errors.New("not implemented")
}

// Post implements [taskutil.Tasker.Post].
func (t tfApply) Post(_ context.Context) error { return nil }
