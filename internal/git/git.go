package igit

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

// Git holds metadata about the current state of the Git repository.
type Git struct {
	Root         string
	Branch       string
	LatestTag    string
	LatestCommit string
	IsDirty      bool
}

// Status holds various pieces of information about Git status.
type Status struct {
	Diff           []string
	UntrackedFiles []string
}

// New returns a populated [Git].
func New(ctx context.Context) (*Git, error) {
	root, err := taskutil.RunCommand(ctx, []string{"git", "rev-parse", "--show-toplevel"})
	if err != nil {
		return nil, err
	}
	iprint.Debugf("Git root on host: '%s'\n", root)

	branch, err := taskutil.RunCommand(ctx, []string{"git", "rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return nil, err
	}
	iprint.Debugf("Git branch: '%s'\n", branch)

	latestTag, err := taskutil.RunCommand(ctx, []string{"bash", "-c", "git tag --list | tail -n1"})
	if err != nil {
		return nil, err
	}
	iprint.Debugf("latest Git tag: '%s'\n", latestTag)

	latestCommit, err := taskutil.RunCommand(ctx, []string{"git", "rev-parse", "--short=8", "HEAD"})
	if err != nil {
		return nil, err
	}
	iprint.Debugf("latest Git commit: '%s'\n", latestCommit)

	gitStatus, err := getRawStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Git status: %w", err)
	}

	var isDirty bool
	if len(gitStatus.Diff) > 0 || len(gitStatus.UntrackedFiles) > 0 {
		isDirty = true
	}

	out := Git{
		Root:         root,
		Branch:       branch,
		LatestTag:    latestTag,
		LatestCommit: latestCommit,
		IsDirty:      isDirty,
	}
	iprint.Debugf("Git: %+v\n", out)

	return &out, nil
}

// SanitizedBranch returns the current branch name, sanitized for various systems that allow for a
// smaller charset (e.g. container image tags).
func (g *Git) SanitizedBranch() string {
	return regexp.MustCompile(`[_/]`).ReplaceAllString(g.Branch, "-")
}

// String implements [fmt.Stringer].
func (g *Git) String() string {
	out := "Current Git information:\n"

	t := reflect.TypeOf(*g)
	v := reflect.ValueOf(*g)
	for i := range v.NumField() {
		field := t.Field(i)
		value := v.Field(i)
		out += fmt.Sprintf("- %s: %v\n", field.Name, value)
	}

	out += "\n"

	return out
}

// getRawStatus returns a slightly-modified "git status" output, so that calling tools can parse it
// more easily.
func getRawStatus(ctx context.Context) (Status, error) {
	outputBytes, err := taskutil.RunCommand(ctx, []string{"git", "status", "--porcelain"})
	if err != nil {
		return Status{}, fmt.Errorf("getting git status output: %w", err)
	}

	output := string(outputBytes)
	outputSplit := strings.Split(output, "\n")

	untrackedFiles := make([]string, 0)
	diff := make([]string, 0)
	for _, line := range outputSplit {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "??") {
			filename := strings.ReplaceAll(line, "?? ", "")
			untrackedFiles = append(untrackedFiles, filename)
		} else {
			filename := regexp.MustCompile(`^( +)?[A-Z]+ +`).ReplaceAllString(line, "")
			diff = append(diff, filename)
		}
	}

	return Status{
		Diff:           diff,
		UntrackedFiles: untrackedFiles,
	}, nil
}
