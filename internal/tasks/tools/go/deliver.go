package gotools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	"github.com/opensourcecorp/oscar/internal/system"
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

	if repo.HasGo {
		out := make([]taskutil.Tasker, 0)

		if cfg.GetDeliverables().GetGoGithubRelease() != nil {
			out = append(out, ghRelease{})
		}

		return out, nil
	}

	return nil, nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t ghRelease) InfoText() string { return "GitHub Release" }

// Exec implements [taskutil.Tasker.Exec].
func (t ghRelease) Exec(ctx context.Context) error {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return err
	}

	var buildErr error
	for _, src := range cfg.GetDeliverables().GetGoGithubRelease().GetBuildSources() {
		buildErr = goBuild(ctx, src)
	}
	if buildErr != nil {
		return err
	}

	buildDir := "build"
	distDir := "dist"

	if err := os.RemoveAll(distDir); err != nil {
		return fmt.Errorf("removing dist directory: %w", err)
	}

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("creating dist directory: %w", err)
	}

	if err := os.CopyFS(distDir, os.DirFS(buildDir)); err != nil {
		return fmt.Errorf("copying build artifacts to %s: %w", distDir, err)
	}

	draftFlag := ""
	if cfg.GetDeliverables().GetGoGithubRelease().GetDraft() {
		draftFlag = "--draft"
	}

	args := []string{"bash", "-c", fmt.Sprintf(`
		gh release create v%s %s --generate-notes --verify-tag --latest ./dist/*
		`, cfg.GetVersion(), draftFlag,
	)}

	if _, err := system.RunCommand(ctx, args); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t ghRelease) Post(_ context.Context) error { return nil }

// goBuild cross-compiles the provided source package and places the resulting artifacts in a
// root-level "build/" subdirectory.
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

		if _, err := system.RunCommand(ctx, []string{"bash", "-c", fmt.Sprintf(`
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
