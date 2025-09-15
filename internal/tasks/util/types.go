package taskutil

import (
	"context"
	"time"
)

// Tasker defines the method set for working with metadata for a given CI Task.
type Tasker interface {
	// InfoText should return a human-readable display string that describes the task, e.g. "Run
	// tests".
	InfoText() string
	// Exec should perform the actual task's actions.
	Exec(ctx context.Context) error
	// Post should perform any post-run actions for the task, if necessary.
	Post(ctx context.Context) error
}

// TaskMap aliases a map of a Task's language/tooling name to its list of Tasks.
type TaskMap map[string][]Tasker

// A Tool defines information about a tool used for running oscar's tasks. A Tool should be defined
// if a language etc. cannot perform the task itself. For example, you would not need a Tool to
// represent a task that runs "go test", but you *would* need a tool to represent a task that runs
// the external "staticcheck" linter for Go.
type Tool struct {
	// TODO
	RunArgs []string
	// The path to the tool's config file, if it has one to use.
	ConfigFilePath string
}

// A Run holds metadata about an instance of an oscar subcommand run, e.g. a list of Tasks for the
// `ci` subcommand.
type Run struct {
	// The "type" of Run as an informative string, i.e. "CI", "Deliver", etc. Used in
	// banner-printing in [Run.PrintRunTypeBanner].
	Type string
	// A timestamp for storing when the overall run started.
	StartTime time.Time
	// Keeps track of all task failures.
	Failures []string
}

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
