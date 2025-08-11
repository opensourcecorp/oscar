package ci

import (
	"errors"
	"fmt"
	"strings"
	"time"

	goci "github.com/opensourcecorp/oscar/internal/ci/go"
	pythonci "github.com/opensourcecorp/oscar/internal/ci/python"
	shellci "github.com/opensourcecorp/oscar/internal/ci/shell"
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	igit "github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// TaskMap is a less-verbose type alias for mapping language names to function signatures that
// return a language's tasks.
type TaskMap map[string][]ciutil.Tasker

// GetCITaskMap assembles the overall list of CI tasks, keyed by their language/tooling name
func GetCITaskMap() (TaskMap, error) {
	repo, err := ciutil.GetRepoComposition()
	if err != nil {
		return nil, fmt.Errorf("getting repo composition: %w", err)
	}

	taskMap := make(TaskMap, 0)
	for langName, getTasksFunc := range map[string]func(ciutil.Repo) []ciutil.Tasker{
		"Go":     goci.Tasks,
		"Python": pythonci.Tasks,
		"Shell":  shellci.Tasks,
		// "Markdown":   markdownci.Tasks,
		// "JavaScript": nodejsci.Tasks,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			taskMap[langName] = tasks
		}
	}

	if len(taskMap) > 0 {
		fmt.Print(repo.String())
		iprint.Debugf("GetCITasks output: %+v\n", taskMap)
	}

	return taskMap, nil
}

// Run defines the behavior for running all CI tasks for the repository.
func Run() (err error) {
	runStartTime := time.Now()

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

	// Handle system init
	if err := ciutil.InitSystem(); err != nil {
		return fmt.Errorf("initializing system: %w", err)
	}

	// Handle all other inits
	fmt.Printf("Initializing the host for discovered file types, this might take some time...\n")
	for _, tasks := range ciTaskMap {
		for _, t := range tasks {
			if err := t.Init(); err != nil {
				return fmt.Errorf("running init for '%s': %w", t.InfoText(), err)
			}
		}
	}
	fmt.Printf("Done with initialization!\n\n")

	// For tracking any changes to Git status etc. after each Task runs
	git, err := igit.New()
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

			// NOTE: this error is checked later
			runErr := t.Run()

			if err := git.Update(); err != nil {
				return fmt.Errorf("internal error: %w", err)
			}
			gitStatusHasChanged, err := git.StatusHasChanged()
			if err != nil {
				return fmt.Errorf("internal error: %w", err)
			}

			if runErr != nil || gitStatusHasChanged {
				iprint.Errorf("FAILED!\n")
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
				git, err = igit.New()
				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}
			} else {
				fmt.Printf("PASSED (t: %s)\n", time.Since(taskStartTime).Round(time.Second/1000).String())
			}
			if err := t.Post(); err != nil {
				iprint.Errorf("running task post-steps: %v\n", err)
				failures = append(failures, fmt.Sprintf("%s :: %s", lang, t.InfoText()))
			}
		}
	}

	if len(failures) > 0 {
		iprint.Errorf("\n================================================================\n")
		iprint.Errorf("The following checks failed and/or caused a git diff:\n")
		for _, f := range failures {
			iprint.Errorf("- %s\n", f)
		}
		iprint.Errorf("================================================================\n\n")
		return errors.New("one or more CI checks failed")
	}

	fmt.Printf("All checks passed! (finished in %s)\n", time.Since(runStartTime).Round(time.Second/1000).String())

	return err
}
