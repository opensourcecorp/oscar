// Package ciconfig is used for storing embeddable config files for various CI tools, that are
// injected at runtime.
package ciconfig

import "embed"

// Files stores config files for each CI tool that will have its own overridden.
//
//go:embed *.conf *.toml
var Files embed.FS
