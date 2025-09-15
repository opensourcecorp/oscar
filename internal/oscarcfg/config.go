package oscarcfg

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"go.yaml.in/yaml/v4"
)

// Config defines the top-level structure of oscar's config file.
type Config struct {
	// Version is the version string for the codebase.
	Version string `yaml:"version" json:"version"`
	// Deliver is the collection of possible deliverable artifacts.
	Deliver *Deliverables `yaml:"deliver" json:"deliver"`
	// Deploy Deployables  `yaml:"deploy" json:"deploy"`
}

// Deliverables contains a field for each possible deliverable.
type Deliverables struct {
	// See [GoGitHubRelease].
	GoGitHubRelease *GoGitHubRelease `yaml:"go_github_release" json:"go_github_release"`
	// See [ContainerImage].
	ContainerImage *ContainerImage `yaml:"container_image" json:"container_image"`
}

// GoGitHubRelease defines the arguments necessary to create GitHub Releases for Go binaries.
type GoGitHubRelease struct {
	// The target GitHub Repository.
	Repo string `yaml:"repo" json:"repo"`
	// The filepaths to the "main" packages to be built.
	BuildSources []string `yaml:"build_sources" json:"build_sources"`
	// Flags whether the Release should be left in Draft state at create-time. This can be useful to
	// set if you want to review the Release contents before actually publishing.
	Draft bool
}

// ContainerImage defines the arguments necessary to build & push container image artifacts.
type ContainerImage struct {
	// The target registry provider domain, e.g. "ghcr.io".
	Registry string `yaml:"registry" json:"registry"`
	// The target OCI repository name, e.g. "oscar".
	Owner string `yaml:"owner" json:"owner"`
	// The target OCI repository, e.g. "oscar".
	Repo string `yaml:"repo" json:"repo"`
}

// Get returns a populated [Config] based on the oscar config file location. If `path` is not
// provided, it will default to looking in the calling directory.
func Get(pathOverride ...string) (Config, error) {
	path := consts.DefaultOscarCfgFileName

	// Handle the override so we can test this function, and use it in other ways (like checking the
	// main branch's version data)
	if len(pathOverride) > 0 {
		path = pathOverride[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading oscar config file: %w", err)
	}
	iprint.Debugf("data read from oscar config file: %s\n", string(data))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshalling oscar config file '%s': %w", path, err)
	}

	return cfg, nil
}
