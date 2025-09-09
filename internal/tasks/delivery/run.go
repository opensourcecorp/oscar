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
	"github.com/opensourcecorp/oscar/internal/tools"
	gotools "github.com/opensourcecorp/oscar/internal/tools/go"
)

// GetDeliveryTaskMap assembles the overall list of Delivery tasks, keyed by their language/tooling
// name.
func GetDeliveryTaskMap(ctx context.Context) (tools.TaskMap, error) {
	repo, err := tools.GetRepoComposition(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting repo composition: %w", err)
	}

	out := make(tools.TaskMap, 0)
	for langName, getTasksFunc := range map[string]func(tools.Repo) []tools.Tasker{
		"Go": gotools.TasksForDelivery,
	} {
		tasks := getTasksFunc(repo)
		if len(tasks) > 0 {
			out[langName] = tasks
		}
	}

	if len(out) > 0 {
		fmt.Print(repo.String())
		iprint.Debugf("GetDeliveryTaskMap output: %#v\n", out)
	}

	return out, nil
}

// Run defines the behavior for running all Delivery tasks for the repository.
func Run(ctx context.Context) (err error) {
	runStartTime := time.Now()

	// Handle system init
	if err := tools.InitSystem(ctx); err != nil {
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

	// All the Delivery tasks that will be looped over. Will also print a summary of discovered file
	// types.
	deliveryTaskMap, err := GetDeliveryTaskMap(ctx)
	if err != nil {
		return fmt.Errorf("getting Delivery tasks: %w", err)
	}

	// Log padding setup
	for lang, tasks := range deliveryTaskMap {
		longestLanguageNameLength = max(longestLanguageNameLength, len(lang))
		for _, t := range tasks {
			longestInfoTextLength = max(longestInfoTextLength, len(t.InfoText()))
		}
	}

	iprint.Debugf("longestLanguageNameLength: %d\n", longestLanguageNameLength)
	iprint.Debugf("longestInfoTextLength: %d\n", longestInfoTextLength)

	// Keeps track of all task failures
	failures := make([]string, 0)
	for lang, tasks := range deliveryTaskMap {
		langNameBannerPadding := strings.Repeat("=", longestLanguageNameLength-len(lang)/2)
		fmt.Printf(
			"%s%s %s %s%s\n",
			strings.Repeat("=", 24), langNameBannerPadding, lang, langNameBannerPadding, strings.Repeat("=", 24),
		)

		for _, t := range tasks {
			// NOTE: if no InfoText() method is provided, it's probably a lang-wide init func, so
			// skip it
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
			runErr = errors.Join(runErr, t.Run(ctx))
			runErr = errors.Join(runErr, t.Post(ctx))

			if runErr != nil {
				iprint.Errorf("FAILED    (%s)\n", tools.RunDurationString(taskStartTime))
				iprint.Errorf("%v\n", runErr)

				failures = append(failures, fmt.Sprintf("%s :: %s", lang, t.InfoText()))
			} else {
				fmt.Printf("SUCCEEDED (%s)\n", tools.RunDurationString(taskStartTime))
			}
		}
	}

	if len(failures) > 0 {
		iprint.Errorf("\n================================================================\n")
		iprint.Errorf("The following tasks failed: (%s)\n", tools.RunDurationString(runStartTime))
		for _, f := range failures {
			iprint.Errorf("- %s\n", f)
		}
		iprint.Errorf("================================================================\n\n")
		return errors.New("one or more Delivery tasks failed")
	}

	fmt.Printf("All tasks succeeded! (%s)\n", tools.RunDurationString(runStartTime))

	return err
}
