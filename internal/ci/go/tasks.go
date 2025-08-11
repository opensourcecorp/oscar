package goci

import (
	"fmt"
	"os"
	"path/filepath"

	ciconfig "github.com/opensourcecorp/oscar/internal/ci/configfiles"
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// A list of tasks that all implement [ciutil.Tasker], for Go.
type (
	baseInitTask     struct{}
	goModCheckTask   struct{}
	goFormatTask     struct{}
	generateCodeTask struct{}
	goBuildTask      struct{}
	goVetTask        struct{}
	staticcheckTask  struct{}
	reviveTask       struct{}
	errcheckTask     struct{}
	goImportsTask    struct{}
	govulncheckTask  struct{}
	goTestTask       struct{}
)

var tasks = []ciutil.Tasker{
	baseInitTask{},
	goModCheckTask{},
	goFormatTask{},
	generateCodeTask{},
	goBuildTask{},
	goVetTask{},
	staticcheckTask{},
	reviveTask{},
	errcheckTask{},
	goImportsTask{},
	govulncheckTask{},
	goTestTask{},
}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasGo {
		return tasks
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t baseInitTask) InfoText() string { return "" }

// Init implements [ciutil.Tasker.Init].
func (t baseInitTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing Go... ")

	goPath := filepath.Join(consts.OscarHome, "go")
	goBinPath := filepath.Join(goPath, "bin")

	// Set the GOPATH
	if err := os.Setenv("GOPATH", goPath); err != nil {
		return fmt.Errorf("setting GOPATH: %w", err)
	}

	// Add the implicit GOBIN to $PATH
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", goBinPath, os.Getenv("PATH"))); err != nil {
		return fmt.Errorf("updating PATH for Go: %w", err)
	}
	iprint.Debugf("PATH after Go init: %s\n", os.Getenv("PATH"))

	if err := os.MkdirAll(goPath, 0755); err != nil {
		return fmt.Errorf("creating GOPATH: %w", err)
	}

	// Now, install Go itself
	if ciutil.IsToolUpToDate(goAsTask) {
		return nil
	}

	hostInput := ciutil.HostInfoInput{
		KernelLinux: "linux",
		KernelMacOS: "darwin",
		ArchAMD64:   "amd64",
		ArchARM64:   "arm64",
	}

	host, err := ciutil.GetHostInfo(hostInput)
	if err != nil {
		return fmt.Errorf("getting host info during init: %w", err)
	}

	releaseURL := fmt.Sprintf(
		goAsTask.RemotePath,
		goAsTask.Version, host.Kernel, host.Arch,
	)

	downloadDir := filepath.Join(os.TempDir(), "go")
	downloadedFile := filepath.Join(downloadDir, "go.tar.gz")

	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
			mkdir -p %s
			curl -fsSL -o %s %s
			rm -rf %s
			tar -C %s -xf %s
		`,
			downloadDir,
			downloadedFile, releaseURL,
			goPath,
			consts.OscarHome, downloadedFile,
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
func (t goModCheckTask) InfoText() string { return "go.mod tidy check" }

// Init implements [ciutil.Tasker.Init].
func (t goModCheckTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t goModCheckTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "mod", "tidy"}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goModCheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goFormatTask) InfoText() string { return "Format" }

// Init implements [ciutil.Tasker.Init].
func (t goFormatTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t goFormatTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "fmt", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goFormatTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t generateCodeTask) InfoText() string { return "Generate code" }

// Init implements [ciutil.Tasker.Init].
func (t generateCodeTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t generateCodeTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "generate", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t generateCodeTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goBuildTask) InfoText() string { return "Build" }

// Init implements [ciutil.Tasker.Init].
func (t goBuildTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t goBuildTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "build", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goBuildTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goVetTask) InfoText() string { return "Vet" }

// Init implements [ciutil.Tasker.Init].
func (t goVetTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t goVetTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "vet", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goVetTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t staticcheckTask) InfoText() string { return "Lint (staticcheck)" }

// Init implements [ciutil.Tasker.Init].
func (t staticcheckTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing staticcheck... ")

	if err := goInstall(staticcheck); err != nil {
		return err
	}

	cfgFileContents, err := ciconfig.Files.ReadFile(filepath.Base(staticcheck.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(staticcheck.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t staticcheckTask) Run() (err error) {
	args := []string{staticcheck.Name, "./..."}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t staticcheckTask) Post() error {
	if err := os.RemoveAll(staticcheck.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t reviveTask) InfoText() string { return "Lint (revive)" }

// Init implements [ciutil.Tasker.Init].
func (t reviveTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing revive... ")

	if err := goInstall(revive); err != nil {
		return err
	}

	cfgFileContents, err := ciconfig.Files.ReadFile(filepath.Base(revive.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(revive.ConfigFilePath, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t reviveTask) Run() error {
	args := []string{
		revive.Name,
		"--config", revive.ConfigFilePath,
		"--set_exit_status",
		"./...",
	}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t reviveTask) Post() error {
	if err := os.RemoveAll(revive.ConfigFilePath); err != nil {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// InfoText implements [ciutil.Tasker.InfoText].
func (t errcheckTask) InfoText() string { return "Lint (errcheck)" }

// Init implements [ciutil.Tasker.Init].
func (t errcheckTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing errcheck... ")

	if err := goInstall(errcheck); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t errcheckTask) Run() error {
	args := []string{errcheck.Name, "./..."}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t errcheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goImportsTask) InfoText() string { return "Format imports" }

// Init implements [ciutil.Tasker.Init].
func (t goImportsTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing goimports... ")

	if err := goInstall(goimports); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t goImportsTask) Run() error {
	args := []string{goimports.Name, "-l", "-w", "."}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goImportsTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t govulncheckTask) InfoText() string { return "Vulnerability scan (govulncheck)" }

// Init implements [ciutil.Tasker.Init].
func (t govulncheckTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Go: Installing govulncheck... ")

	if err := goInstall(govulncheck); err != nil {
		return err
	}

	return nil
}

// Run implements [ciutil.Tasker.Run].
func (t govulncheckTask) Run() error {
	args := []string{govulncheck.Name, "./..."}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t govulncheckTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t goTestTask) InfoText() string { return "Test" }

// Init implements [ciutil.Tasker.Init].
func (t goTestTask) Init() error { return nil }

// Run implements [ciutil.Tasker.Run].
func (t goTestTask) Run() error {
	if err := ciutil.RunCommand([]string{"go", "test", "./..."}); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t goTestTask) Post() error { return nil }
