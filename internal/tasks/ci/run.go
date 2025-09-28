package ci

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
	igit "github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	containertools "github.com/opensourcecorp/oscar/internal/tasks/tools/containers"
	gotools "github.com/opensourcecorp/oscar/internal/tasks/tools/go"
	mdtools "github.com/opensourcecorp/oscar/internal/tasks/tools/markdown"
	pytools "github.com/opensourcecorp/oscar/internal/tasks/tools/python"
	shtools "github.com/opensourcecorp/oscar/internal/tasks/tools/shell"
	versiontools "github.com/opensourcecorp/oscar/internal/tasks/tools/version"
	yamltools "github.com/opensourcecorp/oscar/internal/tasks/tools/yaml"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// getCITaskMap assembles the overall list of CI tasks, keyed by their language/tooling name
func getCITaskMap(repo taskutil.Repo) (taskutil.TaskMap, error) {
	out := make(taskutil.TaskMap)
	for langName, getTasksFunc := range map[string]func(taskutil.Repo) []taskutil.Tasker{
		"Versioning": versiontools.NewTasksForCI,
		"Go":         gotools.NewTasksForCI,
		"Python":     pytools.NewTasksForCI,
		// "Terraform":     tftools.NewTasksForCI,
		"YAML":          yamltools.NewTasksForCI,
		"Containerfile": containertools.NewTasksForCI,
		"Shell":         shtools.NewTasksForCI,
		"Markdown":      mdtools.NewTasksForCI,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	iprint.Debugf("getCITaskMap output: %#v\n", out)

	return out, nil
}

// Run defines the behavior for running all CI tasks for the repository.
func Run(ctx context.Context) (err error) {
	// The mise config that oscar uses is written during init, so be sure to defer its removal here
	defer func() {
		if rmErr := os.RemoveAll(consts.MiseConfigFileName); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing mise config file: %w", rmErr))
		}
	}()

	run, err := taskutil.NewRun(ctx, "CI")
	if err != nil {
		return fmt.Errorf("internal error setting up run info: %w", err)
	}

	taskMap, err := getCITaskMap(run.Repo)
	if err != nil {
		return err
	}

	// For tracking any changes to Git status etc. after each CI Task runs
	gitCI, err := igit.NewForCI(ctx)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
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

			if err := gitCI.Update(ctx); err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			gitStatusHasChanged, err := gitCI.StatusHasChanged(ctx)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}

			if runErr != nil || gitStatusHasChanged {
				iprint.Errorf("FAILED (%s)\n", iprint.RunDurationString(taskStartTime))
				iprint.Errorf("\n")

				if runErr != nil {
					iprint.Errorf("%v\n", runErr)
				}

				if gitStatusHasChanged {
					iprint.Errorf("Files ~CHANGED~ during run: %+v\n", gitCI.CurrentStatus.Diff)
					iprint.Errorf("Files +CREATED+ during run: %+v\n", gitCI.CurrentStatus.UntrackedFiles)
					iprint.Errorf("\n")
				}

				run.Failures = append(run.Failures, fmt.Sprintf("%s :: %s", lang, task.InfoText()))

				// Also need to reset the baseline status
				gitCI, err = igit.NewForCI(ctx)
				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}
			} else {
				iprint.Goodf("PASSED (%s)\n", iprint.RunDurationString(taskStartTime))
			}
		}
	}

	if len(run.Failures) > 0 {
		return run.ReportFailure(err)
	}

	run.ReportSuccess()

	return err
}
