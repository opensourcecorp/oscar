package ci

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run defines the behavior for running all CI tasks for the repository.
func Run() (err error) {
	var (
		// Some sentinel errors to cut down on typing or looking silly elsewhere. Note that since
		// all errors should be logged in this function (and not in their caller, unless during
		// debug), their values aren't actually going to ever be printed anywhere.
		errInternal       = errors.New("internal error")
		errCIChecksFailed = errors.New("one or more CI checks failed")

		// all the CI configs that will be looped over
		ciConfigs = GetCIConfigs()

		// gitDiff string

		// Vars for determining text padding in output banners
		longestLanguageNameLength int
		longestInfoTextLength     int
	)
	for _, c := range ciConfigs {
		longestLanguageNameLength = max(longestLanguageNameLength, len(c.LanguageName))
		for _, t := range c.CITasks {
			longestInfoTextLength = max(longestInfoTextLength, len(t.InfoText))
		}
	}

	// Handle inits
	fmt.Printf("Initializing the host, this might take some time...\n")
	for _, c := range ciConfigs {
		for _, t := range c.CITasks {
			if t.InitScript != "" {
				initSplit := strings.Split(t.InitScript, " ")
				if len(initSplit) <= 1 {
					errPrintf(
						"Internal Error: InitScript struct field '%s:<%s>' was not well-formed\n",
						c.LanguageName, t.InfoText,
					)
					return errInternal
				}

				cmd := exec.Command(initSplit[0], initSplit[1:]...)
				if output, err := cmd.CombinedOutput(); err != nil {
					return fmt.Errorf(
						"running initialization for '%s:<%s>: %w -- output: %s",
						c.LanguageName, t.InfoText, err, string(output),
					)
				}
			}
		}
	}
	fmt.Printf("Done!\n\n")

	failures := make([]string, 0)
	for _, c := range ciConfigs {
		langNameBannerPadding := strings.Repeat("=", longestLanguageNameLength-len(c.LanguageName)/2)
		fmt.Printf(
			"============%s %s %s============\n",
			langNameBannerPadding, c.LanguageName, langNameBannerPadding,
		)

		for _, t := range c.CITasks {
			taskBannerPadding := strings.Repeat(".", longestInfoTextLength-len(t.InfoText))
			// NOTE: no trailing newline on purpose
			fmt.Printf("> %s %s............", t.InfoText, taskBannerPadding)

			splitCmd := strings.Split(t.RunScript, " ")
			cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
			if output, err := cmd.CombinedOutput(); err != nil {
				errPrintf("FAILED!\n\n")
				errPrintf("%s\n", string(output))
				failures = append(failures, t.InfoText)
			} else {
				fmt.Printf("PASSED\n")
			}
		}
	}

	if len(failures) > 0 {
		errPrintf("\n================================================================\n")
		errPrintf("The following checks either failed, or caused a git diff:\n")
		for _, f := range failures {
			errPrintf("- %s\n", f)
		}
		errPrintf("================================================================\n\n")
		return errCIChecksFailed
	}

	fmt.Printf("All checks passed!\n")

	return err
}

// errPrintf is a helper function that writes to standard error.
func errPrintf(format string, args ...any) {
	if _, err := fmt.Fprintf(os.Stderr, format, args...); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(fmt.Sprintf("trying to write to stderr: %v", err))
	}
}
