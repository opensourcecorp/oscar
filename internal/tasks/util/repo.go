package taskutil

import (
	"context"
	"errors"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/system"
)

// A Repo stores information about the contents of the repository being run against.
type Repo struct {
	HasGo            bool
	HasPython        bool
	HasShell         bool
	HasTerraform     bool
	HasContainerfile bool
	HasYaml          bool
	HasMarkdown      bool
}

// String implements the [fmt.Stringer] interface.
func (repo Repo) String() string {
	var out string

	out += "The following file types were found in this repo, and tasks will be run against them:\n"

	if repo.HasGo {
		out += "- Go\n"
	}
	if repo.HasPython {
		out += "- Python\n"
	}
	if repo.HasShell {
		out += "- Shell (sh, bash, etc.)\n"
	}
	if repo.HasTerraform {
		out += "- Terraform\n"
	}
	if repo.HasContainerfile {
		out += "- Containerfile\n"
	}
	if repo.HasYaml {
		out += "- YAML\n"
	}
	if repo.HasMarkdown {
		out += "- Markdown\n"
	}

	// One more newline for padding
	out += "\n"

	return out
}

// NewRepo returns a populated [Repo].
func NewRepo(ctx context.Context) (Repo, error) {
	var errs error

	hasGo, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("go"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasPython, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("py"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasShell, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("sh"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasTerraform, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("tf"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasContainerfile, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("containerfile"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasYaml, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("yaml"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	hasMarkdown, err := system.FilesExistInTree(ctx, system.GetFileTypeListerCommand("md"))
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return Repo{}, errs
	}

	repo := Repo{
		HasGo:            hasGo,
		HasPython:        hasPython,
		HasShell:         hasShell,
		HasTerraform:     hasTerraform,
		HasContainerfile: hasContainerfile,
		HasYaml:          hasYaml,
		HasMarkdown:      hasMarkdown,
	}
	iprint.Debugf("repo composition: %+v\n", repo)

	return repo, nil
}
