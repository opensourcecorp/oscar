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
	// RunScript holds the command that actually runs the task.
	RunScript string
	// InitScript holds an optional command that prepares the host for actually running the command
	// as defined in [CITask.RunScript].
	InitScript string
	// InitScript holds an optional command that runs after the task as defined in
	// [CITask.RunScript].
	PostScript string
}

// GetCIConfigs assembles the overall list of [Config]s and their [Task]s to run.
func GetCIConfigs() []Config {
	return []Config{
		{
			LanguageName: "Go",
			Tasks: []Task{
				{
					InfoText:  "Generate code",
					RunScript: "go generate ./...",
				},
				{
					InfoText:  "Format code",
					RunScript: "go fmt ./...",
				},
				{
					InfoText:  "Build",
					RunScript: "go build ./main.go",
				},
				{
					InfoText:  "Vet",
					RunScript: "go vet ./...",
				},
				{
					InfoText:   "Linter (staticcheck)",
					RunScript:  "go run honnef.co/go/tools/cmd/staticcheck@latest ./...",
					InitScript: "go install honnef.co/go/tools/cmd/staticcheck@latest",
				},
				{
					InfoText:   "Linter (revive)",
					RunScript:  "go run github.com/mgechev/revive@latest --set_exit_status ./...",
					InitScript: "go install github.com/mgechev/revive@latest",
				},
				{
					InfoText:   "Linter (errcheck)",
					RunScript:  "go run github.com/kisielk/errcheck@latest ./...",
					InitScript: "go install github.com/kisielk/errcheck@latest",
				},
				{
					InfoText:   "Vulnerability scanner (govulncheck)",
					RunScript:  "go run golang.org/x/vuln/cmd/govulncheck@latest ./...",
					InitScript: "go install golang.org/x/vuln/cmd/govulncheck@latest",
				},
				{
					InfoText:  "Run tests",
					RunScript: "go test ./...",
				},
			},
		},
	}
}
