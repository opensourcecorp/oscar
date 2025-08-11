package nodejsci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

type (
	// BaseInitTask is exported in this package because there are tools for other languages that
	// provide tooling using Node packages, and so this task can be included in their own tasks to
	// ensure e.g. `npx` availability.
	BaseInitTask struct{}
)

var tasks = []ciutil.Tasker{
	BaseInitTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasNodejs {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t BaseInitTask) InfoText() string { return "" }

// Init implements [ciutil.Tasker.Init].
func (t BaseInitTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- JavaScript: Installing Nodejs... ")

	hostInput := ciutil.HostInfoInput{
		KernelLinux: "linux",
		KernelMacOS: "darwin",
		ArchAMD64:   "x64",
		ArchARM64:   "arm64",
	}

	host, err := ciutil.GetHostInfo(hostInput)
	if err != nil {
		return fmt.Errorf("getting host info during init: %w", err)
	}

	var fileExt string
	if host.Kernel == "linux" {
		fileExt = ".tar.xz"
	} else {
		fileExt = ".tar.gz"
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		nodejs.RemotePath,
		nodejs.Version, nodejs.Version, host.Kernel, host.Arch, fileExt,
	)

	extractedArchiveSubdir := strings.ReplaceAll(filepath.Base(releaseURL), fileExt, "")
	nodeJSHome := filepath.Join(consts.OscarHome, "nodejs", extractedArchiveSubdir)
	cacheDir := filepath.Join(nodeJSHome, "cache")

	downloadDir := filepath.Join(os.TempDir(), "nodejs")
	downloadedFile := filepath.Join(downloadDir, fmt.Sprintf("nodejs%s", fileExt))

	// Add bindir $PATH, once we know the full one. It can be determined via the last part of the
	// release URL.
	binDir := filepath.Join(nodeJSHome, "bin")
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", binDir, os.Getenv("PATH"))); err != nil {
		return fmt.Errorf("updating PATH for Nodejs: %w", err)
	}
	iprint.Debugf("PATH after Nodejs init: %s\n", os.Getenv("PATH"))

	envVars := map[string]string{
		// "npm_config_global":     "true",
		"npm_config_userconfig": filepath.Join(nodeJSHome, ".npmrc"),
		"npm_config_cache":      cacheDir,
		"npm_config_yes":        "true",
	}
	for k, v := range envVars {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("setting env var '%s': %w", k, err)
		}
	}

	// NOW, and ONLY NOW (because we needed so many vars), can we check that npm et al are
	// up-to-date
	if ciutil.IsToolUpToDate(nodejs) {
		return nil
	}

	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
				rm -rf %s
				mkdir -p %s %s %s
				curl -fsSL -o %s %s
				tar -C %s -xf %s
			`,
			nodeJSHome,
			downloadDir, nodeJSHome, cacheDir,
			downloadedFile, releaseURL,
			nodeJSHome, downloadedFile,
		),
	}

	if err := ciutil.RunCommand(installCmd); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t BaseInitTask) Run() error { return nil }

// Post implements [ciutil.Tasker.Post].
func (t BaseInitTask) Post() error { return nil }
