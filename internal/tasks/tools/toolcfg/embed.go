package toolcfg

import "embed"

// Files stores config files for each CI tool.
//
//go:embed *
var Files embed.FS
