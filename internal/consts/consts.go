// Package consts provides a shared place for constants & global variables for usage across the
// codebase.
package consts

import (
	"os"
	"path/filepath"
)

const (
	// DebugEnvVarName is for enabling debug logs etc.
	DebugEnvVarName = "OSC_DEBUG"

	// LangNameGo is the language identifier used for Go configurations.
	LangNameGo = "Go"
	// LangNamePython is the language identifier used for Python configurations.
	LangNamePython = "Python"
	// LangNameSystem is the language identifier used for system configurations.
	LangNameSystem = "System"
)

var (
	// OscarHome is oscar's home directory, under which anything oscar-related will live.
	OscarHome = filepath.Join(os.Getenv("HOME"), ".oscar")
	// OscarHomeBin is the directory where any commands that oscar installs for itself will live.
	OscarHomeBin = filepath.Join(OscarHome, "bin")
)
