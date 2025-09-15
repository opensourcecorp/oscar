package taskutil

import (
	"context"
	"slices"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
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

// RenderRunCommandArgs uses [Tool.RunArgs] and does naive templating to replace certain values
// before being used.
//
// This is useful because when instantiating [Tasker]s, sometimes the fields need to be
// self-referential within the struct. For example, if a [Tool.RunArgs] needs to specify a config
// file path, but it can't know that value until instantiation (even though it's likely defined
// right below it in [Tool.ConfigFilePath]), you can write the `RunArgs` to have a
// `{{ConfigFilePath}}` placeholder that will be interpolated when calling this function.
func (t Tool) RenderRunCommandArgs() []string {
	out := make([]string, len(t.RunArgs))
	for i, arg := range t.RunArgs {
		out[i] = strings.ReplaceAll(arg, `{{ConfigFilePath}}`, t.ConfigFilePath)
	}

	iprint.Debugf("RenderRunCommandArgs: %#v\n", out)

	return out
}

// TaskMap aliases a map of a Task's language/tooling name to its list of Tasks.
type TaskMap map[string][]Tasker

// SortedKeys sorts the keys of the [TaskMap]. Useful for iterating through Tasks in a predictable
// order during runs.
func (tm TaskMap) SortedKeys() []string {
	keys := make([]string, 0)
	for key := range tm {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}
