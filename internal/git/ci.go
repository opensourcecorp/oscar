package git

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tools"
)

// CI defines metadata & behavior for CI tasks.
type CI struct {
	// BaselineStatus is used to check against when running CI checks, so that each CI task can
	// see if it introduced changes.
	BaselineStatus Status
	// CurrentStatus is the latest-available Git status, which may differ from the baseline.
	CurrentStatus Status
}

// Status holds various pieces of information about Git status.
type Status struct {
	Diff           []string
	UntrackedFiles []string
}

// NewForCI returns Git information for CI tasks.
func NewForCI(ctx context.Context) (*CI, error) {
	status, err := getRawStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &CI{
		BaselineStatus: status,
	}, nil
}

// Update recalculates various Git metadata, respecting any existing baseline values set in [NewForCI].
func (g *CI) Update(ctx context.Context) error {
	status, err := getRawStatus(ctx)
	if err != nil {
		return fmt.Errorf("getting Git status: %w", err)
	}

	untrackedFiles := make([]string, 0)
	diff := make([]string, 0)

	for _, line := range status.Diff {
		if !slices.Contains(g.BaselineStatus.Diff, line) {
			filename := regexp.MustCompile(`^ [A-Z] `).ReplaceAllString(line, "")
			diff = append(diff, filename)
		}
	}

	for _, line := range status.UntrackedFiles {
		if !slices.Contains(g.BaselineStatus.UntrackedFiles, line) {
			filename := strings.ReplaceAll(line, "?? ", "")
			untrackedFiles = append(diff, filename)
		}
	}

	g.CurrentStatus = Status{
		Diff:           diff,
		UntrackedFiles: untrackedFiles,
	}

	return nil
}

// StatusHasChanged informs the caller of whether or not the [Status] now differs from the baseline.
func (g *CI) StatusHasChanged(ctx context.Context) (bool, error) {
	if err := g.updateStatus(ctx); err != nil {
		return false, err
	}

	iprint.Debugf("len(g.CurrentStatus.Diff) = %d\n", len(g.CurrentStatus.Diff))
	iprint.Debugf("len(g.BaselineStatusForCI.Diff) = %d\n", len(g.BaselineStatus.Diff))
	iprint.Debugf("len(g.CurrentStatus.UntrackedFiles) = %d\n", len(g.CurrentStatus.UntrackedFiles))
	iprint.Debugf("len(g.BaselineStatusForCI.UntrackedFiles) = %d\n", len(g.BaselineStatus.UntrackedFiles))

	statusChanged := (len(g.CurrentStatus.Diff)+len(g.BaselineStatus.Diff)) != len(g.BaselineStatus.Diff) ||
		(len(g.CurrentStatus.UntrackedFiles)+len(g.BaselineStatus.UntrackedFiles)) != len(g.BaselineStatus.UntrackedFiles)

	iprint.Debugf("statusChanged: %v\n", statusChanged)

	return statusChanged, nil
}

// getRawStatus returns a slightly-modified "git status" output, so that calling tools can parse it
// more easily.
func getRawStatus(ctx context.Context) (Status, error) {
	outputBytes, err := tools.RunCommand(ctx, []string{"git", "status", "--porcelain"})
	if err != nil {
		return Status{}, fmt.Errorf("getting git status output: %w", err)
	}

	output := string(outputBytes)
	outputSplit := strings.Split(output, "\n")

	untrackedFiles := make([]string, 0)
	diff := make([]string, 0)
	for _, line := range outputSplit {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "??") {
			filename := strings.ReplaceAll(line, "?? ", "")
			untrackedFiles = append(untrackedFiles, filename)
		} else {
			filename := regexp.MustCompile(`^( +)?[A-Z]+ +`).ReplaceAllString(line, "")
			diff = append(diff, filename)
		}
	}

	return Status{
		Diff:           diff,
		UntrackedFiles: untrackedFiles,
	}, nil
}

// updateStatus updates the tracked Git status so that it can be compared against the baseline.
func (g *CI) updateStatus(ctx context.Context) error {
	// So any future debug logs have a line break in them
	iprint.Debugf("\n")

	iprint.Debugf("OLD git: %+v\n", g)

	status, err := getRawStatus(ctx)
	if err != nil {
		return fmt.Errorf("getting Git status: %w", err)
	}

	diff := make([]string, 0)
	untrackedFiles := make([]string, 0)

	for _, line := range status.Diff {
		filename := regexp.MustCompile(`^( +)?[A-Z]+ +`).ReplaceAllString(line, "")
		if !slices.Contains(g.BaselineStatus.Diff, filename) {
			diff = append(diff, filename)
		}
	}

	for _, line := range status.UntrackedFiles {
		filename := strings.ReplaceAll(line, "?? ", "")
		if !slices.Contains(g.BaselineStatus.UntrackedFiles, filename) {
			untrackedFiles = append(diff, filename)
		}
	}

	g.CurrentStatus = Status{
		Diff:           diff,
		UntrackedFiles: untrackedFiles,
	}

	iprint.Debugf("NEW git: %+v\n", g)

	return nil
}
