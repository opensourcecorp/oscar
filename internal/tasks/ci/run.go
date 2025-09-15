package ci

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
	"github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	containerfiletools "github.com/opensourcecorp/oscar/internal/tasks/tools/containerfile"
	gotools "github.com/opensourcecorp/oscar/internal/tasks/tools/go"
	mdtools "github.com/opensourcecorp/oscar/internal/tasks/tools/markdown"
	pytools "github.com/opensourcecorp/oscar/internal/tasks/tools/python"
	shtools "github.com/opensourcecorp/oscar/internal/tasks/tools/shell"
	versiontools "github.com/opensourcecorp/oscar/internal/tasks/tools/version"
	yamltools "github.com/opensourcecorp/oscar/internal/tasks/tools/yaml"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// getCITasksMap assembles the overall list of CI tasks, keyed by their language/tooling name
func getCITasksMap(ctx context.Context) (taskutil.TasksMap, error) {
	repo, err := taskutil.GetRepoComposition(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting repo composition: %w", err)
	}

	out := make(taskutil.TasksMap)
	for langName, getTasksFunc := range map[string]func(taskutil.Repo) []taskutil.Tasker{
		"Versioning": versiontools.NewTasksForCI,
		"Go":         gotools.NewTasksForCI,
		"Python":     pytools.NewTasksForCI,
		// "Terraform":     tftools.NewTasksForCI,
		"YAML":          yamltools.NewTasksForCI,
		"Containerfile": containerfiletools.NewTasksForCI,
		"Shell":         shtools.NewTasksForCI,
		"Markdown":      mdtools.NewTasksForCI,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	if len(out) > 0 {
		fmt.Print(repo.String())
		iprint.Debugf("getCITasksMap output: %#v\n", out)
	}

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

	tasksMap, err := getCITasksMap(ctx)
	if err != nil {
		return err
	}

	run, err := taskutil.NewRun(ctx, tasksMap)
	if err != nil {
		return fmt.Errorf("internal error setting up run info: %w", err)
	}

	// For tracking any changes to Git status etc. after each CI Task runs
	gitCI, err := git.NewForCI(ctx)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
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

			if err := gitCI.Update(ctx); err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			gitStatusHasChanged, err := gitCI.StatusHasChanged(ctx)
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}

			if runErr != nil || gitStatusHasChanged {
				iprint.Errorf("FAILED (%s)\n", taskutil.RunDurationString(taskStartTime))
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
				gitCI, err = git.NewForCI(ctx)
				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}
			} else {
				fmt.Printf("PASSED (%s)\n", taskutil.RunDurationString(taskStartTime))
			}
		}
	}

	if len(run.Failures) > 0 {
		iprint.Errorf("\n================================================================\n")
		iprint.Errorf("The following checks failed and/or caused a git diff: (%s)\n", taskutil.RunDurationString(run.StartTime))
		for _, f := range run.Failures {
			iprint.Errorf("- %s\n", f)
		}
		iprint.Errorf("================================================================\n\n")
		return errors.New("one or more CI checks failed")
	}

	fmt.Printf("All checks passed! (%s)\n", taskutil.RunDurationString(run.StartTime))

	return err
}
