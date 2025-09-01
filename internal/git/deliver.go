package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tools"
)

// Delivery defines metadata & behavior for Delivery tasks.
type Delivery struct {
	Root      string
	LatestTag string
	// From {Root}/VERSION file
	CurrentVersionFromFile string
}

// NewForDelivery returns Git information for Delivery tasks.
func NewForDelivery() (*Delivery, error) {
	root, err := tools.RunCommand([]string{"git", "rev-parse", "--show-toplevel"})
	if err != nil {
		return nil, err
	}

	latestTag, err := tools.RunCommand([]string{"bash", "-c", "git tag --list | tail -n1"})
	if err != nil {
		return nil, err
	}

	versionFileContents, err := os.ReadFile(filepath.Join(root, "VERSION"))
	if err != nil {
		return nil, err
	}

	versionFileLines := strings.Split(string(versionFileContents), "\n")
	var version string
	for _, line := range versionFileLines {
		if regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`).MatchString(line) {
			version = line
			break
		}
	}

	if version == "" {
		return nil, fmt.Errorf("could not determine a Semantic Version from your 'VERSION' file")
	}

	out := &Delivery{
		Root:                   root,
		LatestTag:              latestTag,
		CurrentVersionFromFile: version,
	}
	iprint.Debugf("git.Delivery: %+v\n", out)

	return out, nil
}
