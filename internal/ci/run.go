package ci

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	igit "github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// Run defines the behavior for running all CI tasks for the repository.
func Run() (err error) {
	var (
		// Some sentinel errors to cut down on typing or looking silly elsewhere. Note that since
		// all errors should be logged in this function (and not in their caller, unless during
		// debug), their values aren't actually going to ever be printed anywhere.
		errInternal       = errors.New("internal error")
		errCIChecksFailed = errors.New("one or more CI checks failed")

		// All the CI configs that will be looped over
		ciConfigs = GetCIConfigs()

		// Vars for determining text padding in output banners
		longestLanguageNameLength int
		longestInfoTextLength     int
	)

	for _, c := range ciConfigs {
		longestLanguageNameLength = max(longestLanguageNameLength, len(c.LanguageName))
		for _, t := range c.Tasks {
			longestInfoTextLength = max(longestInfoTextLength, len(t.InfoText))
		}
	}

	iprint.Debugf("longestLanguageNameLength: %d\n", longestLanguageNameLength)
	iprint.Debugf("longestInfoTextLength: %d\n", longestInfoTextLength)

	// Handle inits
	fmt.Printf("Initializing the host, this might take some time...\n")
	for _, c := range ciConfigs {
		for _, t := range c.Tasks {
			if t.InitFunc != nil {
				if err := t.InitFunc(); err != nil {
					iprint.Errorf(
						"running initialization for '%s:<%s>: %v",
						c.LanguageName, t.InfoText, err,
					)
					return errInternal
				}
			}
		}
	}
	fmt.Printf("Done!\n\n")

	// For tracking any changes to Git status etc. after each Task runs
	git, err := igit.New()
	if err != nil {
		iprint.Errorf("Internal Error: %v\n", err)
		return errInternal
	}

	failures := make([]string, 0)
	for _, c := range ciConfigs {
		langNameBannerPadding := strings.Repeat("=", longestLanguageNameLength-len(c.LanguageName)/2)
		fmt.Printf(
			"============%s %s %s============\n",
			langNameBannerPadding, c.LanguageName, langNameBannerPadding,
		)

		for _, t := range c.Tasks {
			if t.RunScript == nil {
				continue
			}

			taskBannerPadding := strings.Repeat(".", longestInfoTextLength-len(t.InfoText))
			// NOTE: no trailing newline on purpose
			fmt.Printf("> %s %s............", t.InfoText, taskBannerPadding)

			cmd := exec.Command(t.RunScript[0], t.RunScript[1:]...)

			// NOTE: this error is checked later
			output, runErr := cmd.CombinedOutput()

			if err := git.Update(); err != nil {
				iprint.Errorf("Internal Error: %v\n", err)
				return errInternal
			}
			gitStatusHasChanged, err := git.StatusHasChanged()
			if err != nil {
				iprint.Errorf("Internal Error: %v\n", err)
				return errInternal
			}

			if runErr != nil || gitStatusHasChanged {
				iprint.Errorf("FAILED!\n")
				iprint.Errorf("\n")

				if runErr != nil {
					iprint.Errorf("%s\n", string(output))
				}

				if gitStatusHasChanged {
					iprint.Errorf("Files ~CHANGED~ during run: %+v\n", git.CurrentStatus.Diff)
					iprint.Errorf("Files +CREATED+ during run: %+v\n", git.CurrentStatus.UntrackedFiles)
					iprint.Errorf("\n")
				}

				failures = append(failures, fmt.Sprintf("%s :: %s", c.LanguageName, t.InfoText))

				// Also need to reset the baseline status
				git, err = igit.New()
				if err != nil {
					iprint.Errorf("Internal Error: %v\n", err)
					return errInternal
				}
			} else {
				fmt.Printf("PASSED\n")
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
		return errCIChecksFailed
	}

	fmt.Printf("All checks passed!\n")

	return err
}
