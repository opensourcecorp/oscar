package versiontools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/opensourcecorp/oscar/internal/consts"
	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type versionCI struct{ taskutil.Tool }

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(_ taskutil.Repo) []taskutil.Tasker {
	return []taskutil.Tasker{
		versionCI{},
	}
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t versionCI) InfoText() string { return "Versioning checks" }

// Exec implements [taskutil.Tasker.Exec].
func (t versionCI) Exec(ctx context.Context) (err error) {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return fmt.Errorf("getting oscar config: %w", err)
	}
	version := cfg.GetVersion()
	iprint.Debugf("provided version: %s\n", version)

	// NOTE: we clone the repo in question to a temp location to check the version on
	// the main branch, instead of e.g. trying to checkout main. This is for may reasons, not the
	// least of which being that alternatives would be unreliable in e.g. GitHub Actions CI based on
	// how it treats PR checkouts et al. A small price to pay for reliability.
	tmpCloneDir := filepath.Join(os.TempDir(), "oscar-ci", "this-repo")
	if err := os.MkdirAll(filepath.Dir(tmpCloneDir), 0755); err != nil {
		return fmt.Errorf("creating temp clone parent directory: %w", err)
	}
	defer func() {
		if rmErr := os.RemoveAll(tmpCloneDir); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing temp clone directory: %w", rmErr))
		}
	}()

	remote, err := taskutil.RunCommand(ctx, []string{"git", "remote", "get-url", "origin"})
	if err != nil {
		return fmt.Errorf("determining git root: %w", err)
	}

	remote = canonicalizeGitRemote(remote)

	if _, err := taskutil.RunCommand(ctx, []string{"git", "clone", "--depth", "1", remote, tmpCloneDir}); err != nil {
		return fmt.Errorf("cloning repo source to temp location: %w", err)
	}

	mainCfg, err := oscarcfg.Get(filepath.Join(tmpCloneDir, consts.DefaultOscarCfgFileName))
	if err != nil {
		return fmt.Errorf("getting oscar config: %w", err)
	}
	mainVersion := mainCfg.GetVersion()
	iprint.Debugf("main version: %s\n", version)

	// Need to check if we're already on the main branch, since checking its version against itself
	// will unintentionally fail
	//
	// TODO: update internal git package to have a type with ALL this info so I stop copy-pasting
	// shell-outs around
	branch, err := taskutil.RunCommand(ctx, []string{"git", "rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return fmt.Errorf("checking current Git branch/ref: %w", err)
	}
	iprint.Debugf("current Git branch/ref: %s\n", branch)

	if branch != "main" {
		if !oscarcfg.VersionHasBeenIncremented(version, mainVersion) {
			return fmt.Errorf(
				"version in oscar config on this branch (%s) has not been incremented from the version on the main branch",
				version,
			)
		}
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t versionCI) Post(_ context.Context) error { return nil }

// canonicalizeGitRemote converts a Git remote string to be in canonical HTTPS format.
func canonicalizeGitRemote(remote string) string {
	gitSSHRemoteRegex := regexp.MustCompile(`^(https://|git@)(.*)(:|/)(.*)/(.*(.git)?)$`)
	groups := gitSSHRemoteRegex.FindStringSubmatch(remote)

	out := fmt.Sprintf("https://%s/%s/%s", groups[2], groups[4], groups[5])

	return out
}
