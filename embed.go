// Package oscar sits at the root of the repo, and allows us to embed files all the way down the
// tree.
package oscar

import "embed"

// Files holds any embedded files for use elsewhere across the codebase. Notably, it also holds the
// 'mise.toml' file that is used for not only oscar's own development config but also for its
// internals.
//
//go:embed VERSION mise.toml
var Files embed.FS
