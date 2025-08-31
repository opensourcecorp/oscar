package ci

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tools"
	igo "github.com/opensourcecorp/oscar/internal/tools/go"
	"github.com/opensourcecorp/oscar/internal/tools/markdown"
	"github.com/opensourcecorp/oscar/internal/tools/python"
	"github.com/opensourcecorp/oscar/internal/tools/shell"
)

// TaskMap is a less-verbose type alias for mapping language names to function signatures that
// return a language's tasks.
type TaskMap map[string][]tools.Tasker

// GetCITaskMap assembles the overall list of CI tasks, keyed by their language/tooling name
func GetCITaskMap() (TaskMap, error) {
	repo, err := tools.GetRepoComposition()
	if err != nil {
		return nil, fmt.Errorf("getting repo composition: %w", err)
	}

	out := make(TaskMap, 0)
	for langName, getTasksFunc := range map[string]func(tools.Repo) []tools.Tasker{
		"Go":       igo.Tasks,
		"Python":   python.Tasks,
		"Shell":    shell.Tasks,
		"Markdown": markdown.Tasks,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	if len(out) > 0 {
		fmt.Print(repo.String())
		iprint.Debugf("GetCITasks output: %+v\n", out)
	}

	return out, nil
}

// Run defines the behavior for running all CI tasks for the repository.
func Run() (err error) {
	runStartTime := time.Now()

	// Handle system init
	if err := tools.InitSystem(); err != nil {
		return fmt.Errorf("initializing system: %w", err)
	}
	// The mise config that oscar uses is written during init, so be sure to defer its removal here
	defer func() {
		if rmErr := os.Remove(consts.MiseConfigFileName); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing mise config file: %w", rmErr))
		}
	}()

	var (
		// Vars for determining text padding in output banners
		longestLanguageNameLength int
		longestInfoTextLength     int
	)

	// All the CI tasks that will be looped over. Will also print a summary of discovered file
	// types.
	ciTaskMap, err := GetCITaskMap()
	if err != nil {
		return fmt.Errorf("getting CI tasks: %w", err)
	}

	// Log padding setup
	for lang, tasks := range ciTaskMap {
		longestLanguageNameLength = max(longestLanguageNameLength, len(lang))
		for _, t := range tasks {
			longestInfoTextLength = max(longestInfoTextLength, len(t.InfoText()))
		}
	}

	iprint.Debugf("longestLanguageNameLength: %d\n", longestLanguageNameLength)
	iprint.Debugf("longestInfoTextLength: %d\n", longestInfoTextLength)

	// For tracking any changes to Git status etc. after each Task runs
	git, err := NewGit()
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}

	// Keeps track of all task failures
	failures := make([]string, 0)
	for lang, tasks := range ciTaskMap {
		langNameBannerPadding := strings.Repeat("=", longestLanguageNameLength-len(lang)/2)
		fmt.Printf(
			"============%s %s %s============\n",
			langNameBannerPadding, lang, langNameBannerPadding,
		)

		for _, t := range tasks {
			// NOTE: if no InfoText() method is provided, it's probably a lang-wide init func, so skip it
			if t.InfoText() == "" {
				continue
			}

			taskStartTime := time.Now()

			taskBannerPadding := strings.Repeat(".", longestInfoTextLength-len(t.InfoText()))
			// NOTE: no trailing newline on purpose
			fmt.Printf("> %s %s............", t.InfoText(), taskBannerPadding)

			// NOTE: this error is checked later, when we can check the Run, Post, and git-diff
			// potential errors together
			var runErr error
			runErr = errors.Join(runErr, t.Run())
			runErr = errors.Join(runErr, t.Post())

			if err := git.Update(); err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			gitStatusHasChanged, err := git.StatusHasChanged()
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}

			if runErr != nil || gitStatusHasChanged {
				iprint.Errorf("FAILED! (%s)\n", tools.RunDurationString(taskStartTime))
				iprint.Errorf("\n")

				if runErr != nil {
					iprint.Errorf("%v\n", runErr)
				}

				if gitStatusHasChanged {
					iprint.Errorf("Files ~CHANGED~ during run: %+v\n", git.CurrentStatus.Diff)
					iprint.Errorf("Files +CREATED+ during run: %+v\n", git.CurrentStatus.UntrackedFiles)
					iprint.Errorf("\n")
				}

				failures = append(failures, fmt.Sprintf("%s :: %s", lang, t.InfoText()))

				// Also need to reset the baseline status
				git, err = NewGit()
				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}
			} else {
				fmt.Printf("PASSED (%s)\n", tools.RunDurationString(taskStartTime))
			}
		}
	}

	if len(failures) > 0 {
		iprint.Errorf("\n================================================================\n")
		iprint.Errorf("The following checks failed and/or caused a git diff: (%s)\n", tools.RunDurationString(runStartTime))
		for _, f := range failures {
			iprint.Errorf("- %s\n", f)
		}
		iprint.Errorf("================================================================\n\n")
		return errors.New("one or more CI checks failed")
	}

	fmt.Printf("All checks passed! (%s)\n", tools.RunDurationString(runStartTime))

	return err
}
