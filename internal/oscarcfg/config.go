package oscarcfg

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// Config defines the top-level structure of oscar's config file.
type Config struct {
	// Version is the version string for the codebase.
	Version string `yaml:"version" json:"version"`
	// Deliver is the collection of possible deliverable artifacts.
	Deliver Deliverables `yaml:"deliver" json:"deliver"`
	// Deploy Deployables  `yaml:"deploy" json:"deploy"`
}

// Deliverables contains a field for each possible deliverable.
type Deliverables struct {
	// GoBinaries lists out the Go binaries the user wants to build.
	GoBinaries      []GoBinary      `yaml:"go_binaries" json:"go_binaries"`
	GoGitHubRelease GoGitHubRelease `yaml:"go_github_release" json:"go_github_release"`
}

// GoBinary defines the arguments necessary to build Go binaries. While most other Go-related tasks
// should handle the builds as well, this deliverable type is here to allow users to handle the
// resulting artifacts on their own.
type GoBinary struct {
	// BuildSource is the filepath to the "main" package to be built.
	BuildSource string `yaml:"build_source" json:"build_source"`
}

// GoGitHubRelease defines the arguments necessary to create GitHub Releases for Go binaries.
type GoGitHubRelease struct {
	Repo string `yaml:"repo" json:"repo"`
	// BuildSources are the filepaths to the "main" packages to be built.
	BuildSources []string `yaml:"build_sources" json:"build_sources"`
}

// Get returns a populated [Config] based on the oscar config file location. If `path` is not
// provided, it will default to looking in the calling directory.
func Get(pathOverride ...string) (*Config, error) {
	path := consts.DefaultOscarCfgFileName

	// Handle the override so we can test this function, and use it in other ways (like checking the
	// main branch's version data)
	if len(pathOverride) > 0 {
		path = pathOverride[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading oscar config file: %w", err)
	}
	iprint.Debugf("data read from oscar config file: %s\n", string(data))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling oscar config file '%s': %w", path, err)
	}

	return &cfg, nil
}
