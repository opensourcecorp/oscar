package ci

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
	// InitScript holds an optional command & args that prepares the host for actually running the
	// command as defined in [Task.RunScript].
	InitScript []string
	// InitScript holds an optional command & args that runs after the task as defined in
	// [Task.RunScript].
	PostScript []string
}

// GetCIConfigs assembles the overall list of [Config]s and their [Task]s to run.
func GetCIConfigs() []Config {
	return []Config{
		{
			LanguageName: "Go",
			Tasks: []Task{
				{
					InfoText:  "Go mod check",
					RunScript: []string{"go", "mod", "tidy"},
				},
				{
					InfoText:  "Format code",
					RunScript: []string{"go", "fmt", "./..."},
				},
				{
					InfoText:  "Generate code",
					RunScript: []string{"bash", "-c", "go generate ./... && go fmt ./..."},
				},
				{
					InfoText:  "Build",
					RunScript: []string{"go", "build", "./..."},
				},
				{
					InfoText:  "Vet",
					RunScript: []string{"go", "vet", "./..."},
				},
				{
					InfoText:   "Linter (staticcheck)",
					RunScript:  []string{"go", "run", "honnef.co/go/tools/cmd/staticcheck@latest", "./..."},
					InitScript: []string{"go", "install", "honnef.co/go/tools/cmd/staticcheck@latest"},
				},
				{
					InfoText:   "Linter (revive)",
					RunScript:  []string{"go", "run", "github.com/mgechev/revive@latest", "--set_exit_status", "./..."},
					InitScript: []string{"go", "install", "github.com/mgechev/revive@latest"},
				},
				{
					InfoText:   "Linter (errcheck)",
					RunScript:  []string{"go", "run", "github.com/kisielk/errcheck@latest", "./..."},
					InitScript: []string{"go", "install", "github.com/kisielk/errcheck@latest"},
				},
				{
					InfoText:   "Linter (goimports)",
					RunScript:  []string{"go", "run", "golang.org/x/tools/cmd/goimports@latest", "-l", "-w", "."},
					InitScript: []string{"go", "install", "golang.org/x/tools/cmd/goimports@latest"},
				},
				{
					InfoText:   "Vulnerability scanner (govulncheck)",
					RunScript:  []string{"go", "run", "golang.org/x/vuln/cmd/govulncheck@latest", "./..."},
					InitScript: []string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"},
				},
				{
					InfoText:  "Run tests",
					RunScript: []string{"go", "test", "./..."},
				},
			},
		},
	}
}
