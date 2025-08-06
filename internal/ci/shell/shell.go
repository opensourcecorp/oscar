package shellci

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	"github.com/opensourcecorp/oscar/internal/consts"
)

type (
	shellcheckTask struct{}
	shfmtTask      struct{}
)

var tasks = []ciutil.Tasker{
	shellcheckTask{},
	shfmtTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasShell {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t shellcheckTask) InfoText() string { return "Lint (shellcheck)" }

// Init implements [ciutil.Tasker.Init].
func (t shellcheckTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Shell: Installing shellcheck... ")

	if ciutil.IsCommandUpToDate(shellcheck) {
		return nil
	}

	var shellcheckArch, shellcheckKernel string

	switch runtime.GOARCH {
	case "amd64":
		shellcheckArch = "x86_64"
	case "arm64":
		shellcheckArch = "aarch64"
	default:
		return fmt.Errorf("unsupported CPU architecture '%s'", runtime.GOARCH)
	}

	switch runtime.GOOS {
	case "darwin":
		shellcheckKernel = "darwin"
	case "linux":
		shellcheckKernel = "linux"
	default:
		return fmt.Errorf("unsupported operating system/kernel '%s'", runtime.GOOS)
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		shellcheck.RemotePath,
		shellcheckKernel, shellcheckArch,
	)

	downloadedFile := filepath.Join(os.TempDir(), "shellcheck.tar.xz")

	// NOTE: yes, I know, but this is WAY easier than doing a whole Go song & dance with downloading
	// & unpacking a targz archive. System deps are called out in the README, don't @ me.
	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
					curl -fsSL -o %s %s
					tar -C %s -xf %s
					mv %s/shellcheck-%s/shellcheck %s/
				`,
			downloadedFile, releaseURL,
			os.TempDir(), downloadedFile,
			os.TempDir(), shellcheck.Version, consts.OscarHomeBin,
		),
	}

	if err := ciutil.RunCommand(installCmd); err != nil {
		return err
	}

	return nil
}

func (t shellcheckTask) Run() error {
	args := []string{"bash", "-c", "shopt -s globstar && shellcheck **/*.sh"}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t shellcheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t shfmtTask) InfoText() string { return "Format (shfmt)" }

// Init implements [ciutil.Tasker.Init].
func (t shfmtTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Shell: Installing shfmt... ")

	if ciutil.IsCommandUpToDate(shfmt) {
		return nil
	}

	var shfmtArch, shfmtKernel string

	switch runtime.GOARCH {
	case "amd64":
		shfmtArch = "amd64"
	case "arm64":
		shfmtArch = "arm64"
	default:
		return fmt.Errorf("unsupported CPU architecture '%s'", runtime.GOARCH)
	}

	switch runtime.GOOS {
	case "darwin":
		shfmtKernel = "darwin"
	case "linux":
		shfmtKernel = "linux"
	default:
		return fmt.Errorf("unsupported operating system/kernel '%s'", runtime.GOOS)
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		shfmt.RemotePath,
		shfmt.Version, shfmt.Version, shfmtKernel, shfmtArch,
	)

	downloadedFile := filepath.Join(consts.OscarHomeBin, "shfmt")

	// NOTE: yes, I know, but this is WAY easier than doing a whole Go song & dance with downloading
	// & unpacking a targz archive. System deps are called out in the README, don't @ me.
	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
			curl -fsSL -o %s %s
			chmod +x %s
		`,
			downloadedFile, releaseURL,
			downloadedFile,
		),
	}

	if err := ciutil.RunCommand(installCmd); err != nil {
		return err
	}

	return nil
}

func (t shfmtTask) Run() error {
	args := []string{"bash", "-c", "shopt -s globstar && shfmt **/*.sh"}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t shfmtTask) Post() error { return nil }
