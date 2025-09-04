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

	// MiseVersion is the default version of mise to install if not present. Can be overridden via
	// the `MISE_VERSION` env var, which is checked elsewhere.
	MiseVersion = "v2025.8.21"

	// DefaultOscarCfgFileName is the default basename of oscar's config file.
	DefaultOscarCfgFileName = "oscar.yaml"
)

var (
	// OscarHome is oscar's home directory, under which anything oscar-related will live.
	OscarHome = filepath.Join(os.Getenv("HOME"), ".oscar")
	// OscarHomeBin is the directory where any commands that oscar installs for itself will live.
	OscarHomeBin = filepath.Join(OscarHome, "bin")

	// MiseBinPath is the absolute path to the mise binary, if oscar is the one installing it.
	MiseBinPath = filepath.Join(OscarHomeBin, "mise")

	// MiseConfigFileName is the basename of the mise configuration file that oscar uses
	MiseConfigFileName = "mise.oscar.toml"

	// MiseEnvVars maps mise's env var keys to their desired values.
	MiseEnvVars = map[string]string{
		"MISE_DATA_DIR":  filepath.Join(OscarHome, "share", "mise"),
		"MISE_CACHE_DIR": filepath.Join(OscarHome, "cache", "mise"),
		"MISE_STATE_DIR": filepath.Join(OscarHome, "state", "mise"),
		// used to discover/use mise.<MISE_ENV>.toml, which is the value of [MiseConfigFileName]
		"MISE_ENV": "oscar",
	}
)
