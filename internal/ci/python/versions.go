package pythonci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	runCmd = []string{"uvx"}

	// NOTE: installing will also provide the 'uvx' command, which we also need
	uv = ciutil.VersionedTool{
		Name: "uv",
		// NOTE: this is a pattern string used in macOS & Linux (respectively below) downloads. The
		// positions represent:
		// - version
		// - architecture ("x86_64", or "aarch64")
		// - OS ("apple", or "unknown")
		// - kernel ("darwin", or "linux-gnu")
		RemotePath: "https://github.com/astral-sh/uv/releases/download/%s/uv-%s-%s-%s.tar.gz",
		Version:    "0.8.4",
	}
	ruffLint = ciutil.VersionedTool{
		Name:       "ruff",
		Version:    "0.12.7",
		RunCommand: runCmd,
	}
	ruffFormat = ciutil.VersionedTool{
		Name:       "ruff",
		Version:    "0.12.7",
		RunCommand: runCmd,
	}
	pydoclint = ciutil.VersionedTool{
		Name:       "pydoclint",
		Version:    "0.6.6",
		RunCommand: runCmd,
	}
	mypy = ciutil.VersionedTool{
		Name:       "mypy",
		Version:    "1.17.1",
		RunCommand: runCmd,
	}
)
