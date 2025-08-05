package ciutil

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// InitSystem runs init checks against the host itself, so that other checks can run.
func InitSystem() error {
	requiredSystemCommands := [][]string{
		{"command", "-v", "bash"},
		{"command", "-v", "git"},
		{"command", "-v", "curl"},
		{"command", "-v", "tar"},
	}

	for _, cmd := range requiredSystemCommands {
		iprint.Debugf("Running '%v'\n", cmd)
		if _, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			return fmt.Errorf("command '%s' not found on PATH, cannot continue", cmd[2])
		}
	}

	return nil
}

// RunCommand takes a string slice containing an entire command & its args to run, and returns a
// consisten error message in case of failure.
func RunCommand(cmdArgs []string) error {
	if len(cmdArgs) <= 1 {
		return fmt.Errorf("internal error: not enough arguments passed to RunCommand() -- received: %v", cmdArgs)
	}

	iprint.Debugf("Running '%v'\n", cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(
			"running '%v': %w, with output:\n%s",
			cmdArgs, err, string(output),
		)
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

// func IsCommandVersionOK(t Tasker) bool {
// 	commandVersionOK := true

// 	versionCheckArg := "--version"
// 	if t.CommandVersionCheckArg != "" {
// 		versionCheckArg = t.CommandVersionCheckArg
// 	}

// 	cmd := exec.Command(t.CommandName, versionCheckArg)
// 	outputRaw, err := cmd.CombinedOutput()
// 	output := string(outputRaw)
// 	if err != nil || !strings.Contains(output, t.CommandVersion) {
// 		commandVersionOK = false
// 		if err != nil {
// 			iprint.Debugf("%s possibly not found -- will install\n", t.CommandName)
// 		} else if !strings.Contains(output, t.CommandVersion) {
// 			iprint.Debugf(
// 				"%s version was '%s', but wanted '%s' -- will upgrade for oscar usage\n",
// 				t.CommandName, output, t.CommandVersion,
// 			)
// 		}
// 	}

// 	return commandVersionOK
// }

func GetRepoComposition() (Repo, error) {
	var errs error

	hasGo, err := filesExistInTree("**/*.go")
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasPython, err := filesExistInTree("**/*.py*")
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasShell, err := filesExistInTree("**/*.*sh")
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return Repo{}, errs
	}

	repo := Repo{
		HasGo:     hasGo,
		HasPython: hasPython,
		HasShell:  hasShell,
	}
	iprint.Debugf("repo composition: %+v\n", repo)

	return repo, nil
}
