package ciutil

// Tasker defines the method set for working with metadata for a given CI Task.
type Tasker interface {
	// InfoText should return a human-readable display string that describes the task, e.g. "Run
	// tests". If this is unset, then its banner will not show in the CI log output at all (which
	// may be desirable) in the case of implementers of [Tasker.Init])
	InfoText() string
	// Run should perform the actual task's actions.
	Run() error
	// Post should perform any post-run actions for the task, if necessary.
	Post() error
}

// Repo stores information about the contents of the repository being ran against.
type Repo struct {
	HasGo       bool
	HasPython   bool
	HasShell    bool
	HasNodejs   bool
	HasMarkdown bool
}

// A Tool is a helper struct used to help other types implementing [Tasker] pass around their tool
// versioning/installation information.
type Tool struct {
	// The tool's name, used as an identifier. May also be the tool's invocable command, in which
	// case it can be interpolated as such.
	Name string
	// The installable path for the tool, like a URL. Can also be a format string, e.g. with
	// placeholders for platform-specific strings.
	RemotePath string
	// The version of the tool.
	Version string
	// The path to the tool's config file, if it has one to use.
	ConfigFilePath string
}

// String implements the [fmt.Stringer] interface.
func (repo Repo) String() string {
	var out string

	out += "The following file types were found in this repo, and CI checks will be run against them:\n"

	if repo.HasGo {
		out += "- Go\n"
	}
	if repo.HasPython {
		out += "- Python\n"
	}
	if repo.HasShell {
		out += "- Shell (sh, bash, etc.)\n"
	}
	if repo.HasNodejs {
		out += "- JavaScript (.js/.ts)\n"
	}
	if repo.HasMarkdown {
		out += "- Markdown\n"
	}

	// One more newline for padding
	out += "\n"

	return out
}
