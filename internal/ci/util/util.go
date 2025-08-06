package ciutil

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
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

// IsCommandUpToDate checks whether or not a provided [VersionedTask]'s command is installed, and
// up-to-date on the system. This is used to facilitate skipping unecessary installs, etc.
func IsCommandUpToDate(vt VersionedTask) bool {
	commandUpToDate := true

	version := strings.TrimPrefix(vt.Version, "v")
	re := regexp.MustCompile(version)

	// TODO: revisit how reliable this is for everything
	versionCheckArg := "--version"

	cmd := exec.Command(vt.Name, versionCheckArg)
	outputRaw, err := cmd.CombinedOutput()
	output := string(outputRaw)
	output = strings.Join(strings.Split(output, "\n"), " ")
	if err != nil || !re.MatchString(output) {
		commandUpToDate = false
		if err != nil {
			iprint.Debugf(
				"'%s' possibly not found, or does not have a consistent version flag (error output: %s) -- will install\n",
				vt.Name, output,
			)
		} else if !re.MatchString(output) {
			iprint.Debugf(
				"'%s' version output was '%s', but wanted '%s' -- will upgrade for oscar usage\n",
				vt.Name, output, version,
			)
		}
	}

	if commandUpToDate {
		iprint.Debugf("'%s' found and was up-to-date (version %s)\n", vt.Name, vt.Version)
	}

	return commandUpToDate
}

// GetRepoComposition returns a populated [Repo].
func GetRepoComposition() (Repo, error) {
	var errs error

	hasGo, err := filesExistInTree(`ls **/*.go`)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasPython, err := filesExistInTree(`ls **/*.py*`)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasShell, err := filesExistInTree(`
		# DO NOT discover shell files for oscar itself -- this will make oscar's own CI recurse
		# infinitely
		# TODO: find a way to get around that limitation for oscar's own testing
		if [[ ! $(git remote get-url origin) =~ opensourcecorp/oscar ]] ; then
			find . -type f -name '*.*sh' -or -name '*.bats'
		fi
	`)
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

func filesExistInTree(findScript string) (bool, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		%s
	`, findScript))
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If no files found, that's fine, just report it
		if strings.Contains(string(output), "No such file or directory") {
			return false, nil
		}
		return false, fmt.Errorf("finding files: %w -- output:\n%s", err, string(output))
	}

	if string(output) == "" {
		return false, nil
	}

	return true, nil
}
