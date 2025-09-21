package containertools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	oscarcfgpbv1 "github.com/opensourcecorp/oscar/internal/generated/opensourcecorp/oscar/config/v1"
	igit "github.com/opensourcecorp/oscar/internal/git"
	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
	"go.yaml.in/yaml/v4"
)

type (
	imageBuildPush struct{ taskutil.Tool }
)

// registryMapping contains substructs to be used based on the target OCI registry.
type registryMapping struct {
	GitHub gitHubRegistry
}

// gitHubRegistry provides fields for use in targeting "ghcr.io".
type gitHubRegistry struct {
	// The command to run to authenticate to the registry.
	AuthCommand []string
}

// newRegistryMap returns a populated [registryMapping].
func newRegistryMap(username string) registryMapping {
	return registryMapping{
		GitHub: gitHubRegistry{
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

	if repo.HasContainerfile {
		out := make([]taskutil.Tasker, 0)

		if cfg.GetDeliverables().GetContainerImage() != nil {
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
	cfg := rootCfg.GetDeliverables().GetContainerImage()

	uri, err := constructImageURI(ctx, rootCfg)
	if err != nil {
		return fmt.Errorf("constructing image URI: %w", err)
	}

	composeFileContents, err := os.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}

	composeFile := make(map[string]any)
	if err := yaml.Unmarshal(composeFileContents, composeFile); err != nil {
		return err
	}
	iprint.Debugf("composeFile unmarshalled: %#v\n", composeFile)

	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// GROSS, DUDE
	composeFile["services"].(map[string]any)[cfg.GetName()].(map[string]any)["image"] = uri
	composeFile["services"].(map[string]any)[cfg.GetName()].(map[string]any)["build"].(map[string]any)["context"] = curDir

	composeOut, err := yaml.Marshal(composeFile)
	if err != nil {
		return err
	}
	iprint.Debugf("edited Compose file YAML: %s\n", string(composeOut))

	workDir := filepath.Join(os.TempDir(), "oscar-oci")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}

	outPath := filepath.Join(workDir, "docker-compose.yaml")
	if err := os.WriteFile(outPath, composeOut, 0644); err != nil {
		return err
	}

	registryMap := newRegistryMap(cfg.GetName())

	var authArgs []string
	if strings.Contains(cfg.Registry, "ghcr") {
		authArgs = registryMap.GitHub.AuthCommand
	}

	if _, err := taskutil.RunCommand(ctx, authArgs); err != nil {
		return err
	}

	buildPushArgs := []string{"bash", "-c", fmt.Sprintf(`
		docker compose --file %s build --push %s
		`, outPath, cfg.GetName(),
	)}
	if _, err := taskutil.RunCommand(ctx, buildPushArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t imageBuildPush) Post(_ context.Context) error { return nil }

// constructImageURI constructs an image URI based on data from oscar's config & Git.
func constructImageURI(ctx context.Context, rootCfg *oscarcfgpbv1.Config) (string, error) {
	cfg := rootCfg.GetDeliverables().GetContainerImage()
	git, err := igit.New(ctx)
	if err != nil {
		return "", fmt.Errorf("getting Git info: %w", err)
	}

	tag := rootCfg.GetVersion()
	if git.Branch != "main" {
		tag = fmt.Sprintf("%s-%s", git.SanitizedBranch(), git.LatestCommit)
	}
	if git.IsDirty {
		tag = fmt.Sprintf("%s-%s-dirty", git.SanitizedBranch(), git.LatestCommit)
	}

	uri := fmt.Sprintf(
		"%s/%s/%s:%s",
		cfg.GetRegistry(), cfg.GetNamespace(), cfg.GetName(), tag,
	)
	iprint.Debugf("image URI: %s\n", uri)

	return uri, nil
}
