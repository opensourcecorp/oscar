package ci

import (
	"fmt"
	"os/exec"
	"strings"
)

func runInitCommand(cmdArgs []string) error {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running init: %w -- output:\n%s", err, string(output))
	}

	return nil
}

func filesExistInTree(globstar string) (bool, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("shopt -s globstar && ls %s", globstar))
	if output, err := cmd.CombinedOutput(); err != nil {
		// If no files found, that's fine, just report it
		if strings.Contains(string(output), "No such file or directory") {
			return false, nil
		}
		return false, fmt.Errorf("finding files by globstar: %w -- output:\n%s", err, string(output))
	}

	return true, nil
}
