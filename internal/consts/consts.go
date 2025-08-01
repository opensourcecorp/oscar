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
)

var (
	// OscarHome is oscar's home directory, under which anything oscar-related will live.
	OscarHome = filepath.Join(os.Getenv("HOME"), ".oscar")
	// OscarHomeBin is the directory where any commands that oscar installs for itself will live.
	OscarHomeBin = filepath.Join(OscarHome, "bin")
)

var (
	// GoCIVersions stores Go CI tool versions.
	GoCIVersions = struct {
		Staticcheck string
		Revive      string
		Errcheck    string
		Goimports   string
		Govulncheck string
	}{
		Staticcheck: "2025.1.1",
		Revive:      "v1.11.0",
		Errcheck:    "v1.9.0",
		Goimports:   "v0.35.0",
		Govulncheck: "latest", // yes, on purpose
	}

	// PythonCIVersions stores Python CI tool versions.
	PythonCIVersions = struct {
		UV string
	}{
		UV: "0.8.4",
	}
)
