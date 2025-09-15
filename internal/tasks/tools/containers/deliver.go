package containertools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
	"go.yaml.in/yaml/v4"
)

type (
	imageBuildPush struct{ taskutil.Tool }
)

type registryMapping struct {
	Name   string
	GitHub struct {
		AuthCommand []string
	}
}

func newRegistryMap(username string) registryMapping {
	return registryMapping{
		GitHub: struct{ AuthCommand []string }{
			AuthCommand: []string{"bash", "-c", fmt.Sprintf(`
				echo ${GITHUB_TOKEN} | docker login ghcr.io --username %s --password-stdin
				`, username,
			)},
		},
	}
}

// NewTasksForDelivery returns the list of CI tasks.
func NewTasksForDelivery(repo taskutil.Repo) ([]taskutil.Tasker, error) {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return nil, err
	}

	if repo.HasContainerfile && cfg.Deliver != nil {
		out := make([]taskutil.Tasker, 0)

		if cfg.Deliver.ContainerImage != nil {
			out = append(out, imageBuildPush{})
		}

		return out, nil
	}

	return nil, nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t imageBuildPush) InfoText() string { return "Image Build & Push" }

// Run implements [taskutil.Tasker.Run].
func (t imageBuildPush) Exec(ctx context.Context) error {
	rootCfg, err := oscarcfg.Get()
	if err != nil {
		return err
	}
	cfg := rootCfg.Deliver.ContainerImage

	composeFileContents, err := os.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}

	composeFile := make(map[string]any)
	if err := yaml.Unmarshal(composeFileContents, composeFile); err != nil {
		return err
	}
	iprint.Debugf("composeFile unmarshalled: %#v\n", composeFile)

	// TODO: use a validator package instead, so we can check all the fields more easily
	if cfg.Repo == "" {
		return fmt.Errorf("required 'repo' key not set for this Deliverable in oscar.yaml")
	}

	uri := fmt.Sprintf(
		"%s/%s/%s:%s",
		cfg.Registry, cfg.Owner, cfg.Repo, rootCfg.Version,
	)

	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// GROSS, DUDE
	composeFile["services"].(map[string]any)[cfg.Repo].(map[string]any)["image"] = uri
	composeFile["services"].(map[string]any)[cfg.Repo].(map[string]any)["build"].(map[string]any)["context"] = curDir

	composeOut, err := yaml.Marshal(composeFile)
	if err != nil {
		return err
	}

	workDir := filepath.Join(os.TempDir(), "oscar-oci")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}

	outPath := filepath.Join(workDir, "docker-compose.yaml")
	if err := os.WriteFile(outPath, composeOut, 0644); err != nil {
		return err
	}

	registryMap := newRegistryMap(cfg.Repo)

	var authArgs []string
	if strings.Contains(cfg.Registry, "ghcr") {
		authArgs = registryMap.GitHub.AuthCommand
	}

	if _, err := taskutil.RunCommand(ctx, authArgs); err != nil {
		return err
	}

	buildPushArgs := []string{"bash", "-c", fmt.Sprintf(`
		docker compose --file %s build --push %s
		`, outPath, cfg.Repo,
	)}
	if _, err := taskutil.RunCommand(ctx, buildPushArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t imageBuildPush) Post(_ context.Context) error { return nil }
