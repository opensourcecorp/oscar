package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	ghRelease struct{ taskutil.Tool }
)

// NewTasksForDelivery returns the list of Delivery tasks.
func NewTasksForDelivery(repo taskutil.Repo) ([]taskutil.Tasker, error) {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return nil, err
	}

	if repo.HasGo && cfg.Deliver != nil {
		out := make([]taskutil.Tasker, 0)

		if cfg.Deliver.GoGitHubRelease != nil {
			out = append(out, ghRelease{})
		}

		return out, nil
	}

	return nil, nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t ghRelease) InfoText() string { return "GitHub Releases" }

// Exec implements [taskutil.Tasker.Exec].
func (t ghRelease) Exec(ctx context.Context) error {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return err
	}

	var buildErr error
	for _, src := range cfg.Deliver.GoGitHubRelease.BuildSources {
		buildErr = goBuild(ctx, src)
	}
	if buildErr != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ghRelease) Post(_ context.Context) error { return nil }

func goBuild(ctx context.Context, src string) error {
	if strings.HasSuffix(src, ".go") {
		return fmt.Errorf("provided Go build source '%s' was a file, but must be a path to a package", src)
	}

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

		binName := filepath.Base(src)
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
