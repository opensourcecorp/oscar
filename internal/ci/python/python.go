package pythonci

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	"github.com/opensourcecorp/oscar/internal/consts"
)

type (
	baseInitTask   struct{}
	buildTask      struct{}
	ruffLintTask   struct{}
	ruffFormatTask struct{}
	pydoclintTask  struct{}
	mypyTask       struct{}
)

var tasks = []ciutil.Tasker{
	baseInitTask{},
	buildTask{},
	ruffLintTask{},
	ruffFormatTask{},
	pydoclintTask{},
	mypyTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasPython {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t baseInitTask) InfoText() string { return "" }

// Init implements [ciutil.Tasker.Init].
func (t baseInitTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Python: Installing uv... ")

	if ciutil.IsCommandUpToDate(uv) {
		return nil
	}

	var uvArch, uvOS, uvKernel string

	switch runtime.GOARCH {
	case "amd64":
		uvArch = "x86_64"
	case "arm64":
		uvArch = "aarch64"
	default:
		return fmt.Errorf("unsupported CPU architecture '%s'", runtime.GOARCH)
	}

	switch runtime.GOOS {
	case "darwin":
		uvOS = "apple"
		uvKernel = "darwin"
	case "linux":
		uvOS = "unknown"
		uvKernel = "linux-gnu"
	default:
		return fmt.Errorf("unsupported operating system '%s'", runtime.GOOS)
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		uv.RemotePath,
		uv.Version, uvArch, uvOS, uvKernel,
	)

	// Grab the last element of the download URL (minus the extension) to get the unpacked archive
	// directory name
	urlSplit := strings.Split(releaseURL, "/")
	archiveName := strings.ReplaceAll(urlSplit[len(urlSplit)-1], ".tar.gz", "")

	downloadedFile := filepath.Join(os.TempDir(), "uv.tar.gz")

	// NOTE: yes, I know, but this is WAY easier than doing a whole Go song & dance with downloading
	// & unpacking a targz archive. System deps are called out in the README, don't @ me.
	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
					curl -fsSL -o %s %s
					tar -C %s -xzf %s
					mv %s/%s/{uv,uvx} %s/
				`,
			downloadedFile, releaseURL,
			os.TempDir(), downloadedFile,
			os.TempDir(), archiveName, consts.OscarHomeBin,
		),
	}

	if err := ciutil.RunCommand(installCmd); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t baseInitTask) Run() error { return nil }

// Post implements [ciutil.Tasker.Post].
func (t baseInitTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t buildTask) InfoText() string { return "Build" }

// Init implements [ciutil.Tasker.Init].
func (t buildTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t buildTask) Run() error {
	if err := ciutil.RunCommand([]string{"uv", "build"}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t buildTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t ruffLintTask) InfoText() string { return "Lint (ruff)" }

// Init implements [ciutil.Tasker.Init].
func (t ruffLintTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t ruffLintTask) Run() error {
	args := slices.Concat(getVersionedArgs(ruff), []string{"check", "--fix", "./src"})
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t ruffLintTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t ruffFormatTask) InfoText() string { return "Format (ruff)" }

// Init implements [ciutil.Tasker.Init].
func (t ruffFormatTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t ruffFormatTask) Run() error {
	args := slices.Concat(getVersionedArgs(ruff), []string{"format", "./src"})
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t ruffFormatTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t pydoclintTask) InfoText() string { return "Lint (pydoclint)" }

// Init implements [ciutil.Tasker.Init].
func (t pydoclintTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t pydoclintTask) Run() error {
	args := slices.Concat(getVersionedArgs(pydoclint), []string{"./src"})
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t pydoclintTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t mypyTask) InfoText() string { return "Type-check (mypy)" }

// Init implements [ciutil.Tasker.Init].
func (t mypyTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t mypyTask) Run() error {
	args := slices.Concat(getVersionedArgs(mypy), []string{"./src"})
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t mypyTask) Post() error { return nil }
