package igit

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// Git defines metadata & behavior for Git interactions.
type Git struct {
	// BaselineStatusForCI is used to check against when running CI checks, so that each CI task can
	// see if it introduced changes.
	BaselineStatusForCI Status
	CurrentStatus       Status
}

// Status holds various pieces of information about Git status.
type Status struct {
	Diff           []string
	UntrackedFiles []string
}

// New returns a snapshot of Git information available at call-time.
func New() (*Git, error) {
	status, err := getRawStatus()
	if err != nil {
		return nil, err
	}

	return &Git{
		BaselineStatusForCI: status,
	}, nil
}

// Update recalculates various Git metadata, respecting any existing baseline values set in [New].
func (g *Git) Update() error {
	status, err := getRawStatus()
	if err != nil {
		return err
	}

	untrackedFiles := make([]string, 0)
	diff := make([]string, 0)

	for _, line := range status.Diff {
		if !slices.Contains(g.BaselineStatusForCI.Diff, line) {
			filename := regexp.MustCompile(`^ [A-Z] `).ReplaceAllString(line, "")
			diff = append(diff, filename)
		}
	}

	for _, line := range status.UntrackedFiles {
		if !slices.Contains(g.BaselineStatusForCI.UntrackedFiles, line) {
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
func (g *Git) StatusHasChanged() (bool, error) {
	if err := g.updateStatus(); err != nil {
		return false, err
	}

	iprint.Debugf("len(g.CurrentStatus.Diff) = %d\n", len(g.CurrentStatus.Diff))
	iprint.Debugf("len(g.BaselineStatusForCI.Diff) = %d\n", len(g.BaselineStatusForCI.Diff))
	iprint.Debugf("len(g.CurrentStatus.UntrackedFiles) = %d\n", len(g.CurrentStatus.UntrackedFiles))
	iprint.Debugf("len(g.BaselineStatusForCI.UntrackedFiles) = %d\n", len(g.BaselineStatusForCI.UntrackedFiles))

	statusChanged := (len(g.CurrentStatus.Diff)+len(g.BaselineStatusForCI.Diff)) != len(g.BaselineStatusForCI.Diff) ||
		(len(g.CurrentStatus.UntrackedFiles)+len(g.BaselineStatusForCI.UntrackedFiles)) != len(g.BaselineStatusForCI.UntrackedFiles)

	iprint.Debugf("statusChanged: %v\n", statusChanged)

	return statusChanged, nil
}

func getRawStatus() (Status, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	outputBytes, err := cmd.CombinedOutput()
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

func (g *Git) updateStatus() error {
	// So any future debug logs have a line break in them
	iprint.Debugf("\n")

	iprint.Debugf("OLD git: %+v\n", g)

	status, err := getRawStatus()
	if err != nil {
		return err
	}

	diff := make([]string, 0)
	untrackedFiles := make([]string, 0)

	for _, line := range status.Diff {
		filename := regexp.MustCompile(`^( +)?[A-Z]+ +`).ReplaceAllString(line, "")
		if !slices.Contains(g.BaselineStatusForCI.Diff, filename) {
			diff = append(diff, filename)
		}
	}

	for _, line := range status.UntrackedFiles {
		filename := strings.ReplaceAll(line, "?? ", "")
		if !slices.Contains(g.BaselineStatusForCI.UntrackedFiles, filename) {
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
