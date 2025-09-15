package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	ghRelease struct{ taskutil.Tool }
)

// TasksForDelivery returns the list of Delivery tasks.
func TasksForDelivery(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasGo {
		return []taskutil.Tasker{
			ghRelease{
				Tool: taskutil.Tool{},
			},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t ghRelease) InfoText() string { return "GitHub Release" }

// Exec implements [taskutil.Tasker.Exec].
func (t ghRelease) Exec(ctx context.Context) error {
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

		if _, err := taskutil.RunCommand(ctx, []string{"bash", "-c", fmt.Sprintf(`
			CGO_ENABLED=0 \
			GOOS=%s GOARCH=%s \
			go build -ldflags '-extldflags "-static"' -o %s %s`,
			goos, goarch,
			target, src,
		)}); err != nil {
			return fmt.Errorf("building Go binary: %w", err)
		}

		if err := os.Chmod(target, 0755); err != nil {
			return fmt.Errorf("marking target as executable: %w", err)
		}
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ghRelease) Post(_ context.Context) error { return nil }
