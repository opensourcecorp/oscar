package igit

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	BaselineStatus Status
	CurrentStatus  Status
}

type Status struct {
	Diff           []string
	UntrackedFiles []string
}

func New() (*Git, error) {
	status, err := getStatus()
	if err != nil {
		return nil, err
	}

	return &Git{
		BaselineStatus: status,
	}, nil
}

func (g *Git) Update() error {
	status, err := getStatus()
	if err != nil {
		return err
	}

	g.CurrentStatus = status
	return nil
}

func (g *Git) HasChanged() bool {
	return len(g.CurrentStatus.Diff) != len(g.BaselineStatus.Diff) ||
		len(g.CurrentStatus.UntrackedFiles) != len(g.BaselineStatus.UntrackedFiles)
}

func getStatus() (Status, error) {
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
		if strings.HasPrefix(line, "??") {
			filename := strings.ReplaceAll(line, "?? ", "")
			untrackedFiles = append(untrackedFiles, filename)
		} else {
			// Getting here means that there were existing, diffed files found in the status
			diff = append(diff, line)
		}
	}

	return Status{
		Diff:           diff,
		UntrackedFiles: untrackedFiles,
	}, nil
}
