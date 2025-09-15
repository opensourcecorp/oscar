package taskutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// A Run holds metadata about an instance of an oscar subcommand run, e.g. a list of Tasks for the
// `ci` subcommand.
type Run struct {
	// The "type" of Run as an informative string, i.e. "CI", "Deliver", etc. Used in
	// banner-printing in [Run.PrintRunTypeBanner].
	Type string
	// A timestamp for storing when the overall run started.
	StartTime time.Time
	// Keeps track of all task failures.
	Failures []string
}

// NewRun returns a populated [Run].
func NewRun(ctx context.Context, runType string) (Run, error) {
	// Kind of wonky, but print the banner first
	Run{Type: runType}.PrintRunTypeBanner()

	// Handle system init
	if err := InitSystem(ctx); err != nil {
		return Run{}, fmt.Errorf("initializing system: %w", err)
	}

	return Run{
		Type:      runType,
		StartTime: time.Now(),
		Failures:  make([]string, 0),
	}, nil
}

// PrintRunTypeBanner prints a banner about the type of [Run] underway.
func (run Run) PrintRunTypeBanner() {
	// this padding accounts for leading text length before the run.Type string
	padding := 9
	bannerChar := "#"
	fmt.Printf("%s\n", strings.Repeat(bannerChar, len(run.Type)+padding))
	fmt.Printf("%s Run: %s #\n", bannerChar, run.Type)
	fmt.Printf("%s\n\n", strings.Repeat(bannerChar, len(run.Type)+padding))
}

// PrintTaskMapBanner prints a banner about the [TaskMap] being run.
func (run Run) PrintTaskMapBanner(lang string) {
	fmt.Printf(
		"=== %s %s>\n",
		lang, strings.Repeat("=", 64-len(lang)),
	)
}

// PrintTaskBanner prints a banner about the Task being run.
func (run Run) PrintTaskBanner(task Tasker) {
	// NOTE: no trailing newline on purpose
	fmt.Printf("> %s %s............", task.InfoText(), strings.Repeat(".", 32-len(task.InfoText())))
}

// ReportSuccess prints information about the success of a [Run].
func (run Run) ReportSuccess() {
	fmt.Printf("\nAll tasks succeeded! (%s)\n\n", RunDurationString(run.StartTime))
}

// ReportFailure prints information about the failure of a [Run]. It takes an `error` arg in case
// the caller is expecting to return a joined error because of e.g. deferred calls or later-checked
// errors that an outer variable already holds.
func (run Run) ReportFailure(err error) error {
	iprint.Errorf("\n%s\n", strings.Repeat("=", 65))
	iprint.Errorf("The following tasks failed: (%s)\n", RunDurationString(run.StartTime))
	for _, f := range run.Failures {
		iprint.Errorf("- %s\n", f)
	}
	iprint.Errorf("%s\n\n", strings.Repeat("=", 65))

	err = errors.Join(err, errors.New("one or more tasks failed"))
	return err
}
