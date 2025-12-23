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
	"github.com/opensourcecorp/oscar/internal/system"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
	"go.yaml.in/yaml/v4"
)

type (
	imageBuildPush struct{ taskutil.Tool }

	// RegistryMapping contains substructs to be used based on the target OCI registry.
	registryMapping struct {
		GitHub gitHubRegistry
	}

	// GitHubRegistry provides fields for use in targeting "ghcr.io".
	gitHubRegistry struct {
		// The command to run to authenticate to the registry.
		AuthCommand []string
	}
)

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
		return nil, fmt.Errorf("getting oscarcfg: %w", err)
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

// Exec implements [taskutil.Tasker.Exec].
func (t imageBuildPush) Exec(ctx context.Context) error {
	rootCfg, cfgErr := oscarcfg.Get()
	if cfgErr != nil {
		return fmt.Errorf("getting oscarcfg: %w", cfgErr)
	}

	cfg := rootCfg.GetDeliverables().GetContainerImage()

	uri, uriErr := constructImageURI(ctx, rootCfg)
	if uriErr != nil {
		return fmt.Errorf("constructing image URI: %w", uriErr)
	}

	// TODO: replace all of the below with the actual Compose API usage, e.g. as found in the README
	// example here: https://github.com/compose-spec/compose-go

	composeFileContents, readErr := os.ReadFile("docker-compose.yaml")
	if readErr != nil {
		return fmt.Errorf("reading docker-compose.yaml file: %w", readErr)
	}

	composeFile := make(map[string]any)
	if err := yaml.Unmarshal(composeFileContents, composeFile); err != nil {
		return fmt.Errorf("unmarshalling YAML: %w", err)
	}

	iprint.Debugf("composeFile unmarshalled: %#v\n", composeFile)

	curDir, wdErr := os.Getwd()
	if wdErr != nil {
		return fmt.Errorf("getting workdir: %w", wdErr)
	}

	// TODO: GROSS, DUDE. See comment above about using the actual Compose API.
	composeFile["services"].(map[string]any)[cfg.GetName()].(map[string]any)["image"] = uri
	composeFile["services"].(map[string]any)[cfg.GetName()].(map[string]any)["build"].(map[string]any)["context"] = curDir

	composeOut, mErr := yaml.Marshal(composeFile)
	if mErr != nil {
		return fmt.Errorf("marshalling YAML: %w", mErr)
	}

	iprint.Debugf("edited Compose file YAML: %s\n", string(composeOut))

	workDir := filepath.Join(os.TempDir(), "oscar-oci")
	if err := os.MkdirAll(workDir, 0700); err != nil {
		return fmt.Errorf("making Compose directory: %w", err)
	}

	outPath := filepath.Join(workDir, "docker-compose.yaml")
	if err := os.WriteFile(outPath, composeOut, 0600); err != nil {
		return fmt.Errorf("writing Compose file contents: %w", err)
	}

	registryMap := newRegistryMap(cfg.GetName())

	var authArgs []string
	if strings.Contains(cfg.GetRegistry(), "ghcr") {
		authArgs = registryMap.GitHub.AuthCommand
	}

	if _, err := system.RunCommand(ctx, authArgs); err != nil {
		return fmt.Errorf("running registry auth command: %w", err)
	}

	buildPushArgs := []string{"bash", "-c", fmt.Sprintf(`
		docker compose --file %s build --push %s
		`, outPath, cfg.GetName(),
	)}
	if _, err := system.RunCommand(ctx, buildPushArgs); err != nil {
		return fmt.Errorf("running image build & push command: %w", err)
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

	if git.Branch == "main" {
		// NOTE: this might be a risky hack, but because the containerized tests copy other language
		// files to the root directory and rename them, the git diff will be nonzero, and the tag
		// value will be impossible to check against. So, override that here.
		git.IsDirty = false
	} else {
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
