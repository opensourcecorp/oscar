// Package oscar sits at the root of the repo, and allows us to embed files all the way down the
// tree.
package oscar

import "embed"

// Files holds any embedded files for use elsewhere across the codebase.
//
//go:embed VERSION
var Files embed.FS
