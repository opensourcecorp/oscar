package ciutil

// Tasker defines the method set for working with metadata for a given CI Task.
type Tasker interface {
	// InfoText should return a human-readable display string that describes the task, e.g. "Run
	// tests". If this is unset, then its banner will not show in the CI log output at all (which
	// may be desirable) in the case of implementers of [Tasker.Init])
	InfoText() string
	// Init should perform any task-level initialization for the implementing language. This method
	// can also be used on a "dummy" task to perform language-wide initialization, if desired.
	Init() error
	// Run should perform the actual task's actions.
	Run() error
	Post() error
}

// Repo stores information about the contents of the repository being ran against.
type Repo struct {
	HasGo     bool
	HasPython bool
	HasShell  bool
}

// A VersionedTask is a helper struct used to help other types implementing [Tasker] pass around
// their versioning/installation information.
type VersionedTask struct {
	Name           string
	RemotePath     string
	Version        string
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

	// One more newline for padding
	out += "\n"

	return out
}
