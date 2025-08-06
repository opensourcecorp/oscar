package shellci

import (
	"fmt"
	"os"
	"path/filepath"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
	"github.com/opensourcecorp/oscar/internal/consts"
)

type (
	shellcheckTask struct{}
	shfmtTask      struct{}
	batsTask       struct{}
)

var tasks = []ciutil.Tasker{
	shellcheckTask{},
	shfmtTask{},
	batsTask{},
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

	hostInput := ciutil.HostInfoInput{
		KernelLinux: "linux",
		KernelMacOS: "darwin",
		ArchAMD64:   "x86_64",
		ArchARM64:   "aarch64",
	}

	host, err := ciutil.GetHostInfo(hostInput)
	if err != nil {
		return fmt.Errorf("getting host info during init: %w", err)
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		shellcheck.RemotePath,
		host.Kernel, host.Arch,
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
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shellcheck.Name)}
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

	hostInput := ciutil.HostInfoInput{
		KernelLinux: "linux-gnu",
		KernelMacOS: "darwin",
		ArchAMD64:   "amd64",
		ArchARM64:   "arm64",
	}

	host, err := ciutil.GetHostInfo(hostInput)
	if err != nil {
		return fmt.Errorf("getting host info during init: %w", err)
	}

	// This will also be the name of the directory once extracted from the archive
	releaseURL := fmt.Sprintf(
		shfmt.RemotePath,
		shfmt.Version, shfmt.Version, host.Kernel, host.Arch,
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
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		ls **/*.sh || exit 0
		%s **/*.sh
	`, shfmt.Name)}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t shfmtTask) Post() error { return nil }

// InfoText implements [ciutil.Tasker.InfoText].
func (t batsTask) InfoText() string { return "Test (bats)" }

// Init implements [ciutil.Tasker.Init].
func (t batsTask) Init() error {
	defer fmt.Println("Done.")
	fmt.Printf("- Shell: Installing bats... ")

	if ciutil.IsCommandUpToDate(bats) {
		return nil
	}

	clonePath := filepath.Join(os.TempDir(), "bats")

	// NOTE: yes, I know, but this is WAY easier than doing a whole Go song & dance with downloading
	// & unpacking a targz archive. System deps are called out in the README, don't @ me.
	installCmd := []string{"bash", "-c",
		fmt.Sprintf(`
			git clone %s %s
			git -C %s checkout %s
			bash %s/install.sh %s
		`,
			bats.RemotePath, clonePath,
			clonePath, bats.Version,
			clonePath, consts.OscarHome,
		),
	}

	if err := ciutil.RunCommand(installCmd); err != nil {
		return err
	}

	return nil
}

func (t batsTask) Run() error {
	args := []string{"bash", "-c", fmt.Sprintf(`
		shopt -s globstar
		# Don't run if no bats files found, otherwise it will error out
		ls **/*.bats || exit 0
		%s **/*.bats
	`, bats.Name)}
	if err := ciutil.RunCommand(args); err != nil {
		return err
	}

	return nil
}

// Post implements [ciutil.Tasker.Post].
func (t batsTask) Post() error { return nil }
