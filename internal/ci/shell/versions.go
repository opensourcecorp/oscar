package shellci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var (
	shellcheck = ciutil.VersionedTask{
		Name: "shellcheck",
		// Placeholders are for:
		// - kernel ("darwin", or "linux")
		// - arch ("x86_64", or "aarch")
		RemotePath: "https://github.com/koalaman/shellcheck/releases/download/stable/shellcheck-stable.%s.%s.tar.xz",
		Version:    "v0.11.0",
	}
	shfmt = ciutil.VersionedTask{
		Name: "shfmt",
		// Placeholders are for:
		// - version
		// - version (again)
		// - kernel ("darwin", or "linux")
		// - arch ("amd64", or "arm64")
		RemotePath: "https://github.com/mvdan/sh/releases/download/%s/shfmt_%s_%s_%s",
		Version:    "v3.12.0",
	}
)
