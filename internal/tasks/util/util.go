package taskutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/opensourcecorp/oscar"
	"github.com/opensourcecorp/oscar/internal/consts"
	"github.com/opensourcecorp/oscar/internal/hostinfo"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// InitSystem runs setup & checks against the host itself, so that oscar can run.
func InitSystem(ctx context.Context) error {
	fmt.Printf("Initializing the host, this might take some time... ")
	startTime := time.Now()

	requiredSystemCommands := [][]string{
		{"bash", "--version"},
		{"git", "--version"},
	}

	for _, cmd := range requiredSystemCommands {
		iprint.Debugf("Running '%v'\n", cmd)
		if output, err := exec.CommandContext(ctx, cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			return fmt.Errorf(
				"command '%s' possibly not found on PATH, cannot continue (error: %w -- output: %s)",
				cmd[0], err, string(output),
			)
		}
	}

	for _, d := range []string{consts.OscarHome, consts.OscarHomeBin} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf(
				"internal error when creating oscar directory '%s': %v",
				d, err,
			)
		}
	}

	for name, value := range consts.MiseEnvVars {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf(
				"internal error when setting mise env var '%s': %v",
				name, err,
			)
		}
	}

	if err := installMise(ctx); err != nil {
		return fmt.Errorf("installing mise: %w", err)
	}

	cfgFileContents, err := oscar.Files.ReadFile("mise.toml")
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(consts.MiseConfigFileName, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	// Init for task runs
	if _, err := RunCommand(ctx, []string{consts.MiseBinPath, "trust", consts.MiseConfigFileName}); err != nil {
		return fmt.Errorf("running mise trust: %w", err)
	}
	if _, err := RunCommand(ctx, []string{consts.MiseBinPath, "install"}); err != nil {
		return fmt.Errorf("running mise install: %w", err)
	}

	fmt.Printf("Done! (%s)\n\n", RunDurationString(startTime))

	return nil
}

// NewRun returns a populated [Run].
func NewRun(ctx context.Context, runType string) (Run, error) {
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
	padding := 9
	fmt.Printf("%s\n", strings.Repeat("@", len(run.Type)+padding))
	fmt.Printf("@ Run: %s @\n", run.Type)
	fmt.Printf("%s\n\n", strings.Repeat("@", len(run.Type)+padding))
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
	fmt.Printf("All tasks succeeded! (%s)\n\n", RunDurationString(run.StartTime))
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

// RunCommand takes a string slice containing a command & its args to run, and returns a consistent
// error message in case of failure. It also returns the command output, in case the caller needs to
// parse it on their own.
func RunCommand(ctx context.Context, cmdArgs []string) (string, error) {
	if len(cmdArgs) <= 1 {
		return "", fmt.Errorf("internal error: not enough arguments passed to RunCommand() -- received: %v", cmdArgs)
	}

	var args []string
	if cmdArgs[0] == consts.MiseBinPath {
		args = cmdArgs[1:]
	} else {
		args = slices.Concat([]string{"exec", "--"}, cmdArgs)
	}

	cmd := exec.CommandContext(ctx, consts.MiseBinPath, args...)
	iprint.Debugf("Running '%v'\n", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"running '%v': %w, with output:\n%s",
			cmd.Args, err, string(output),
		)
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

// GetRepoComposition returns a populated [Repo].
func GetRepoComposition(ctx context.Context) (Repo, error) {
	var errs error

	hasGo, err := filesExistInTree(ctx, GetFileTypeListerCommand("go"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasPython, err := filesExistInTree(ctx, GetFileTypeListerCommand("py"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasShell, err := filesExistInTree(ctx, GetFileTypeListerCommand("sh"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasTerraform, err := filesExistInTree(ctx, GetFileTypeListerCommand("tf"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasContainerfile, err := filesExistInTree(ctx, GetFileTypeListerCommand("containerfile"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasYaml, err := filesExistInTree(ctx, GetFileTypeListerCommand("yaml"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasMarkdown, err := filesExistInTree(ctx, GetFileTypeListerCommand("md"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return Repo{}, errs
	}

	repo := Repo{
		HasGo:            hasGo,
		HasPython:        hasPython,
		HasShell:         hasShell,
		HasTerraform:     hasTerraform,
		HasContainerfile: hasContainerfile,
		HasYaml:          hasYaml,
		HasMarkdown:      hasMarkdown,
	}
	iprint.Debugf("repo composition: %+v\n", repo)

	return repo, nil
}

// GetFileTypeListerCommand takes a [ripgrep]-known file type, and returns a slice containing a
// command & its args to find matching files. ripgrep is used as the default file-finder across the
// codebase because not only is it fast, but it also supports files like `.gitignore` without any
// extra configuration.
//
// [ripgrep]: https://github.com/BurntSushi/ripgrep
func GetFileTypeListerCommand(fileType string) string {
	// NOTE: there are some special cases we need to handle, like how ripgrep understands "docker"
	// as a file type arg (and it will find files matching the glob "*Dockerfile*"), but it will
	// *not* find e.g. "Containerfile"
	switch fileType {
	case "containerfile":
		return `rg --files --glob-case-insensitive --glob='*{Containerfile,Dockerfile}*' || true`
	default:
		return fmt.Sprintf(`rg --files --type '%s' || true`, fileType)
	}
}

// RunDurationString returns a calculated duration used to indicate how long a particular Task (or
// set of Tasks) took to run.
func RunDurationString(t time.Time) string {
	return fmt.Sprintf("t: %s", time.Since(t).Round(time.Second/1000).String())
}

// RenderRunCommandArgs uses [Tool.RunArgs] and does naive templating to replace certain values
// before being used.
//
// This is useful because when instantiating [Tasker]s, sometimes the fields need to be
// self-referential within the struct. For example, if a [Tool.RunArgs] needs to specify a config
// file path, but it can't know that value until instantiation (even though it's likely defined
// right below it in [Tool.ConfigFilePath]), you can write the `RunArgs` to have a
// `{{ConfigFilePath}}` placeholder that will be interpolated when calling this function.
func (t Tool) RenderRunCommandArgs() []string {
	out := make([]string, len(t.RunArgs))
	for i, arg := range t.RunArgs {
		out[i] = strings.ReplaceAll(arg, `{{ConfigFilePath}}`, t.ConfigFilePath)
	}

	iprint.Debugf("RenderRunCommandArgs: %#v\n", out)

	return out
}

// installMise installs [mise] into [consts.OscarHomeBin], if not found there.
//
// [mise]: https://mise.jdx.dev
func installMise(_ context.Context) (err error) {
	miseFound := true
	_, err = os.Stat(consts.MiseBinPath)
	if err != nil {
		iprint.Debugf("error when running os.Stat(consts.MiseBinPath): %w\n", err)
		if os.IsNotExist(err) {
			miseFound = false
			iprint.Debugf("mise not found, will install\n")
		} else {
			return fmt.Errorf("internal error checking if mise is installed: %w", err)
		}
	}

	if miseFound {
		// TODO: mise version check
		iprint.Debugf("mise found, nothing to do\n")
		return
	}

	miseVersion := os.Getenv("MISE_VERSION")
	if miseVersion == "" {
		miseVersion = consts.MiseVersion
	}

	hostInput := hostinfo.Input{
		KernelLinux: "linux",
		KernelMacOS: "macos",
		ArchAMD64:   "x64",
		ArchARM64:   "arm64",
	}

	host, err := hostinfo.Get(hostInput)
	if err != nil {
		return fmt.Errorf("getting host info during mise install: %w", err)
	}

	miseReleaseURL := fmt.Sprintf(
		"https://github.com/jdx/mise/releases/download/%s/mise-%s-%s-%s",
		consts.MiseVersion, consts.MiseVersion, host.Kernel, host.Arch,
	)

	out, err := os.Create(consts.MiseBinPath)
	if err != nil {
		return fmt.Errorf("creating mise target file: %w", err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing mise target file: %w", closeErr))
		}
	}()

	// TODO: use a context func instead
	resp, err := http.Get(miseReleaseURL)
	if err != nil {
		return fmt.Errorf("making GET request for mise GitHub Release: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status code when getting mise: %s", resp.Status)
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("writing mise data to target: %w", err)
	}

	if err := os.Chmod(consts.MiseBinPath, 0755); err != nil {
		return fmt.Errorf("changing mise binary to be executable: %w", err)
	}

	return err
}

// filesExistInTree performs file discovery by allowing various tools to check if they need to run
// based on file presence.
func filesExistInTree(ctx context.Context, findScript string) (bool, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		%s`,
		findScript,
	))
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
