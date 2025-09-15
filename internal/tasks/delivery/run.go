package delivery

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	gotools "github.com/opensourcecorp/oscar/internal/tasks/tools/go"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// getDeliveryTasksMap assembles the overall list of Delivery tasks, keyed by their language/tooling
// name.
func getDeliveryTasksMap(ctx context.Context) (taskutil.TasksMap, error) {
	repo, err := taskutil.GetRepoComposition(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting repo composition: %w", err)
	}

	out := make(taskutil.TasksMap)
	for langName, getTasksFunc := range map[string]func(taskutil.Repo) []taskutil.Tasker{
		// "Version": versiontools.TasksForDelivery,
		"Go": gotools.TasksForDelivery,
		// "Python":   pytools.TasksForDelivery,
		// "YAML": yamltools.TasksForDelivery,
		// "Shell":    shtools.TasksForDelivery,
		// "Markdown": mdtools.TasksForDelivery,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	if len(out) > 0 {
		fmt.Print(repo.String())
		iprint.Debugf("getDeliveryTasksMap output: %#v\n", out)
	}

	return out, nil
}

// Run defines the behavior for running all Delivery tasks for the repository.
func Run(ctx context.Context) (err error) {
	// The mise config that oscar uses is written during init, so be sure to defer its removal here
	defer func() {
		if rmErr := os.RemoveAll(consts.MiseConfigFileName); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing mise config file: %w", rmErr))
		}
	}()

	tasksMap, err := getDeliveryTasksMap(ctx)
	if err != nil {
		return err
	}

	run, err := taskutil.NewRun(ctx, tasksMap)
	if err != nil {
		return fmt.Errorf("internal error setting up run info: %w", err)
	}

	for lang, tasks := range tasksMap {
		langNameBannerPadding := strings.Repeat("=", run.LongestLanguageNameLength-len(lang)/2)
		fmt.Printf(
			"%s%s %s %s%s\n",
			strings.Repeat("=", 24), langNameBannerPadding, lang, langNameBannerPadding, strings.Repeat("=", 24),
		)

		for _, task := range tasks {
			taskStartTime := time.Now()

			taskBannerPadding := strings.Repeat(".", run.LongestInfoTextLength-len(task.InfoText()))
			// NOTE: no trailing newline on purpose
			fmt.Printf("> %s %s............", task.InfoText(), taskBannerPadding)

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
		iprint.Errorf("\n================================================================\n")
		iprint.Errorf("The following tasks failed: (%s)\n", taskutil.RunDurationString(run.StartTime))
		for _, f := range run.Failures {
			iprint.Errorf("- %s\n", f)
		}
		iprint.Errorf("================================================================\n\n")
		return errors.New("one or more Delivery tasks failed")
	}

	fmt.Printf("All tasks succeeded! (%s)\n", taskutil.RunDurationString(run.StartTime))

	return err
}
