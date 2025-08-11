package ciutil

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// InitSystem runs init checks against the host itself, so that other checks can run.
func InitSystem() error {
	requiredSystemCommands := [][]string{
		{"bash", "--version"},
		{"git", "--version"},
		{"curl", "--version"},
		{"tar", "--version"},
	}

	for _, cmd := range requiredSystemCommands {
		iprint.Debugf("Running '%v'\n", cmd)
		if output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			return fmt.Errorf(
				"command '%s' possibly not found on PATH, cannot continue (error: %w -- output: %s)",
				cmd[0], err, string(output),
			)
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

// RunVersionedCommand wraps the [RunCommand] function in this package by passing args for the tool
// to be run by a different program, like `uvx` or `npx`.
func RunVersionedCommand(t VersionedTool, args []string) error {
	allArgs := slices.Concat(
		t.RunCommand,
		[]string{fmt.Sprintf("%s@%s", t.Name, t.Version)},
		args,
	)

	if err := RunCommand(allArgs); err != nil {
		return err
	}

	return nil
}

// IsToolUpToDate checks whether or not a provided [VersionedTool]'s tooling is installed, and
// up-to-date on the system. This is used to facilitate skipping unecessary installs, etc.
func IsToolUpToDate(vt VersionedTool) bool {
	commandUpToDate := true

	version := strings.TrimPrefix(vt.Version, "v")
	re := regexp.MustCompile(version)

	versionCheckArg := "--version"
	if vt.VersionCheckArg != "" {
		versionCheckArg = vt.VersionCheckArg
	}

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

	hasPython, err := filesExistInTree(`find . -type f -name '*.py' -or -name '*.pyi'`)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasShell, err := filesExistInTree(`find . -type f -name '*.*sh' -or -name '*.bats'`)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasNodejs, err := filesExistInTree(`find . -type f -name '*.js' -or -name '*.ts'`)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasMarkdown, err := filesExistInTree(`ls **/*.md`)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return Repo{}, errs
	}

	repo := Repo{
		HasGo:       hasGo,
		HasPython:   hasPython,
		HasShell:    hasShell,
		HasNodejs:   hasNodejs,
		HasMarkdown: hasMarkdown,
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
