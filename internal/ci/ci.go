package ci

import (
	"errors"
	"fmt"
	"os"

	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// Config holds metadata for a given language, and its individual CI tasks.
type Config struct {
	// LanguageName is the human-readable name of the programming language having [Task]s run
	// against it.
	LanguageName string
	// Tasks holds the [Task]s to be run.
	Tasks []Task
}

// Task holds metadata for a given CI task.
type Task struct {
	// InfoText is a human-readable display string that describes the task, e.g. "Run tests".
	InfoText string
	// RunScript holds the command & args that actually runs the task.
	RunScript []string
	// InitFunc holds an optional command & args that prepares the host for actually running the
	// command as defined in [Task.RunScript].
	InitFunc func() error
	// InitScript holds an optional command & args that runs after the task as defined in
	// [Task.RunScript].
	PostFunc func() error
}

// Repo stores information about the contents of the repository being ran against.
type Repo struct {
	HasGo     bool
	HasPython bool
	HasShell  bool
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

// GetCIConfigs assembles the overall list of [Config]s and their [Task]s to run.
func GetCIConfigs() []Config {
	repo, err := getRepoComposition()
	if err != nil {
		iprint.Errorf("getting repo composition: %v\n", err)
		os.Exit(1)
	}

	configs := make([]Config, 0)
	for _, funcName := range []func(Repo) []Config{
		getGoConfigs,
		getPythonConfigs,
		getShellConfigs,
	} {
		configs = append(configs, funcName(repo)...)
	}

	if len(configs) > 0 {
		fmt.Print(repo.String())
	}

	return configs
}

func getRepoComposition() (Repo, error) {
	var errs error

	hasGo, err := filesExistInTree("**/*.go")
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasPython, err := filesExistInTree("**/*.py*")
	if err != nil {
		errs = errors.Join(errs, err)
	}
	hasShell, err := filesExistInTree("**/*.*sh")
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return Repo{}, errs
	}

	repo := Repo{
		HasGo:     hasGo,
		HasPython: hasPython,
		HasShell:  hasShell,
	}
	iprint.Debugf("repo composition: %+v\n", repo)

	return repo, nil
}
