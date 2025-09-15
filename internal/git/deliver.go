package git

import (
	"context"
	"fmt"

	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// Delivery defines metadata & behavior for Delivery tasks.
type Delivery struct {
	Root      string
	LatestTag string
	// From oscar config file
	CurrentVersion string
}

// NewForDelivery returns Git information for Delivery tasks.
func NewForDelivery(ctx context.Context) (*Delivery, error) {
	root, err := taskutil.RunCommand(ctx, []string{"git", "rev-parse", "--show-toplevel"})
	if err != nil {
		return nil, err
	}

	latestTag, err := taskutil.RunCommand(ctx, []string{"bash", "-c", "git tag --list | tail -n1"})
	if err != nil {
		return nil, err
	}
	iprint.Debugf("latest Git tag: '%s'\n", latestTag)

	cfg, err := oscarcfg.Get()
	if err != nil {
		return nil, fmt.Errorf("getting oscar config: %w", err)
	}
	version := cfg.Version

	if version == "" {
		return nil, fmt.Errorf("could not determine a Semantic Version from your oscar config file")
	}

	out := &Delivery{
		Root:           root,
		LatestTag:      latestTag,
		CurrentVersion: version,
	}
	iprint.Debugf("git.Delivery: %+v\n", out)

	return out, nil
}
