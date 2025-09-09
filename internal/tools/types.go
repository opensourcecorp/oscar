package tools

import "context"

// Tasker defines the method set for working with metadata for a given CI Task.
type Tasker interface {
	// InfoText should return a human-readable display string that describes the task, e.g. "Run
	// tests". If this is unset, then its banner will not show in the CI log output at all (which
	// may be desirable) in the case of implementers of [Tasker.Init])
	InfoText() string
	// Run should perform the actual task's actions.
	Run(ctx context.Context) error
	// Post should perform any post-run actions for the task, if necessary.
	Post(ctx context.Context) error
}

// Repo stores information about the contents of the repository being ran against.
type Repo struct {
	HasGo        bool
	HasPython    bool
	HasShell     bool
	HasTerraform bool
	HasYaml      bool
	HasMarkdown  bool
}

// TaskMap is a less-verbose type alias for mapping language names to function signatures that
// return a language's tasks.
type TaskMap map[string][]Tasker

// A Tool defines information about a tool used for running oscar's tasks. A Tool should be defined
// if a language etc. cannot perform the task itself. For example, you would not need a Tool to
// represent a task that runs "go test", but you *would* need a tool to represent a task that runs
// the external "staticcheck" linter for Go.
//
// Every Tool should implement [Tasker].
type Tool struct {
	// The tool's name, used as an identifier. May also be the tool's invocable command, in which
	// case it can be interpolated as such.
	Name string
	// The path to the tool's config file, if it has one to use.
	ConfigFilePath string
	// The optional installable path for the tool, like a URL. Can also be a format string, e.g.
	// with placeholders for platform-specific strings. Should mostly not be needed if using mise.
	RemotePath string
	// The version of the tool. Should mostly not be needed if using mise.
	Version string
}

// String implements the [fmt.Stringer] interface.
func (repo Repo) String() string {
	var out string

	out += "The following file types were found in this repo, and tasks can/will be run against them:\n"

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
