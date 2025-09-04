package igo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/oscar/internal/tools"
)

type (
	ghRelease struct{}
)

var deliveryTasks = []tools.Tasker{
	ghRelease{},
}

// TasksForDelivery returns the list of Delivery tasks.
func TasksForDelivery(repo tools.Repo) []tools.Tasker {
	if repo.HasGo {
		return deliveryTasks
	}

	return nil
}

// InfoText implements [tools.Tasker.InfoText].
func (t ghRelease) InfoText() string { return "GitHub Release" }

// Run implements [tools.Tasker.Run].
func (t ghRelease) Run() error {
	targetDir := "build"

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("removing build directory: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating build directory: %w", err)
	}

	distros := []string{
		"linux/amd64",
		"linux/arm64",
		"darwin/amd64",
		"darwin/arm64",
	}

	for _, distro := range distros {
		splits := strings.Split(distro, "/")
		goos := splits[0]
		goarch := splits[1]

		binName := "oscar"

		src := "./cmd/oscar"
		target := filepath.Join(targetDir, fmt.Sprintf("%s-%s-%s", binName, goos, goarch))

		if _, err := tools.RunCommand([]string{"bash", "-c", fmt.Sprintf(`
			GOOS=%s GOARCH=%s go build -o %s %s`,
			goos, goarch, target, src,
		)}); err != nil {
			return fmt.Errorf("building Go binary: %w", err)
		}

		if err := os.Chmod(target, 0755); err != nil {
			return fmt.Errorf("marking target as executable: %w", err)
		}
	}

	return nil
}

// Post implements [tools.Tasker.Post].
func (t ghRelease) Post() error { return nil }
