package delivery

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
	igit "github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tasks/ci"
	containertools "github.com/opensourcecorp/oscar/internal/tasks/tools/containers"
	gotools "github.com/opensourcecorp/oscar/internal/tasks/tools/go"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// getDeliveryTaskMap assembles the overall list of Delivery tasks, keyed by their language/tooling
// name.
func getDeliveryTaskMap(repo taskutil.Repo) (taskutil.TaskMap, error) {
	out := make(taskutil.TaskMap)
	for langName, getTasksFunc := range map[string]func(taskutil.Repo) ([]taskutil.Tasker, error){
		"Go":         gotools.NewTasksForDelivery,
		"OCI Images": containertools.NewTasksForDelivery,
		// "Python":     pytools.NewTasksForDelivery,
		// "Terraform":     tftools.NewTasksForDelivery,
		// "Markdown":      mdtools.NewTasksForDelivery,
	} {
		tasks, err := getTasksFunc(repo)
		if err != nil {
			return nil, err
		}

		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	iprint.Debugf("getDeliveryTaskMap output: %#v\n", out)

	return out, nil
}

// Run defines the behavior for running all Delivery tasks for the repository.
func Run(ctx context.Context) (err error) {
	// We intentionally run CI tasks before allowing any Delivery tasks to begin
	if err := ci.Run(ctx); err != nil {
		return fmt.Errorf("running CI tasks before Delivery tasks: %w", err)
	}

	// The mise config that oscar uses is written during init, so be sure to defer its removal here
	defer func() {
		if rmErr := os.RemoveAll(consts.MiseConfigFileName); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing mise config file: %w", rmErr))
		}
	}()

	run, err := taskutil.NewRun(ctx, "Deliver")
	if err != nil {
		return fmt.Errorf("internal error setting up run info: %w", err)
	}

	git, err := igit.New(ctx)
	if err != nil {
		return err
	}
	fmt.Print(git.String())

	repo, err := taskutil.NewRepo(ctx)
	if err != nil {
		return fmt.Errorf("getting repo composition: %w", err)
	}
	fmt.Print(repo.String())

	taskMap, err := getDeliveryTaskMap(repo)
	if err != nil {
		return err
	}

	for _, lang := range taskMap.SortedKeys() {
		tasks := taskMap[lang]

		run.PrintTaskMapBanner(lang)

		for _, task := range tasks {
			taskStartTime := time.Now()
			run.PrintTaskBanner(task)

			// NOTE: this error is checked later, when we can check the Run, Post, and git-diff
			// potential errors together
			var runErr error
			runErr = errors.Join(runErr, task.Exec(ctx))
			runErr = errors.Join(runErr, task.Post(ctx))

			if runErr != nil {
				iprint.Errorf("FAILED    (%s)\n", taskutil.RunDurationString(taskStartTime))
				iprint.Errorf("%v\n", runErr)

				run.Failures = append(run.Failures, fmt.Sprintf("%s :: %s", lang, task.InfoText()))
			} else {
				fmt.Printf("SUCCEEDED (%s)\n", taskutil.RunDurationString(taskStartTime))
			}
		}
	}

	if len(run.Failures) > 0 {
		return run.ReportFailure(err)
	}

	run.ReportSuccess()

	return err
}
