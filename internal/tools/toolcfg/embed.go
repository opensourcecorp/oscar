// Package toolcfg is used for storing embeddable config files for various tools, that are injected
// at runtime.
package toolcfg

import "embed"

// Files stores config files for each CI tool.
//
//go:embed *.conf *.toml *.yaml
var Files embed.FS
