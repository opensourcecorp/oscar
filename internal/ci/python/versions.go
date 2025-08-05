package pythonci

import (
	"fmt"

	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	// NOTE: installing will also provide the 'uvx' command, which we also need
	uv = ciutil.VersionedTask{
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
	ruff = ciutil.VersionedTask{
		Name:    "ruff",
		Version: "0.12.7",
	}
	pydoclint = ciutil.VersionedTask{
		Name:    "pydoclint",
		Version: "0.6.6",
	}
	mypy = ciutil.VersionedTask{
		Name:    "mypy",
		Version: "1.17.1",
	}
)

func getVersionedArgs(t ciutil.VersionedTask) []string {
	return []string{"uvx", fmt.Sprintf("%s@%s", t.Name, t.Version)}
}
