package taskutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	igit "github.com/opensourcecorp/oscar/internal/git"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/system"
)

// A Run holds metadata about an instance of an oscar subcommand run, e.g. a list of Tasks for the
// `ci` subcommand.
type Run struct {
	// The "type" of Run as an informative string, i.e. "CI", "Deliver", etc. Used in
	// banner-printing in [Run.PrintRunTypeBanner].
	Type string
	// See [igit.Git].
	Git *igit.Git
	// See [Repo].
	Repo Repo
	// See [iprint.AllColors].
	Colors iprint.AllColors
	// A timestamp for storing when the overall run started.
	StartTime time.Time
	// Keeps track of all task failures.
	Failures []string
}

// NewRun returns a populated [Run].
func NewRun(ctx context.Context, runType string) (Run, error) {
	// Kind of wonky, but print the banner first
	Run{Type: runType}.PrintRunTypeBanner()

	colors := iprint.Colors()

	// Handle system init
	if err := system.Init(ctx); err != nil {
		return Run{}, fmt.Errorf("initializing system: %w", err)
	}

	git, err := igit.New(ctx)
	if err != nil {
		return Run{}, err
	}
	iprint.Infof(colors.Gray + git.String() + colors.Reset)

	repo, err := NewRepo(ctx)
	if err != nil {
		return Run{}, fmt.Errorf("getting repo composition: %w", err)
	}
	iprint.Infof(colors.Gray + repo.String() + colors.Reset)

	return Run{
		Type:      runType,
		Git:       git,
		Repo:      repo,
		Colors:    colors,
		StartTime: time.Now(),
		Failures:  make([]string, 0),
	}, nil
}

// PrintRunTypeBanner prints a banner about the type of [Run] underway.
func (run Run) PrintRunTypeBanner() {
	colors := iprint.Colors()

	// this padding accounts for leading text length before the run.Type string
	padding := 9
	bannerChar := "#"
	iprint.Infof(
		colors.Yellow+"%s\n"+colors.Reset,
		strings.Repeat(bannerChar, len(run.Type)+padding),
	)
	iprint.Infof(
		colors.Yellow+"%s Run: %s #\n"+colors.Reset,
		bannerChar, run.Type,
	)
	iprint.Infof(
		colors.Yellow+"%s\n\n"+colors.Reset,
		strings.Repeat(bannerChar, len(run.Type)+padding),
	)
}

// PrintTaskMapBanner prints a banner about the [TaskMap] being run.
func (run Run) PrintTaskMapBanner(lang string) {
	iprint.Infof("=== %s %s>\n", lang, strings.Repeat("=", 64-len(lang)))
}

// PrintTaskBanner prints a banner about the Task being run.
func (run Run) PrintTaskBanner(task Tasker) {
	// NOTE: no trailing newline on purpose
	iprint.Infof(
		"> %s %s............",
		iprint.Colors().White+task.InfoText()+iprint.Colors().InfoColor,
		strings.Repeat(".", 32-len(task.InfoText())),
	)
}

// ReportSuccess prints information about the success of a [Run].
func (run Run) ReportSuccess() {
	iprint.Goodf("\nAll tasks succeeded! (%s)\n\n", iprint.RunDurationString(run.StartTime))
}

// ReportFailure prints information about the failure of a [Run]. It takes an `error` arg in case
// the caller is expecting to return a joined error because of e.g. deferred calls or later-checked
// errors that an outer variable already holds.
func (run Run) ReportFailure(err error) error {
	iprint.Errorf("\n%s\n", strings.Repeat("=", 65))
	iprint.Errorf("The following tasks failed: (%s)\n", iprint.RunDurationString(run.StartTime))
	for _, f := range run.Failures {
		iprint.Errorf("- %s\n", f)
	}
	iprint.Errorf("%s\n\n", strings.Repeat("=", 65))

	err = errors.Join(err, errors.New("one or more tasks failed"))
	return err
}
