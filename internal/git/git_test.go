package igit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBranchForURI(t *testing.T) {
	want := "feature-wip"

	tt := []struct {
		Name   string
		Branch string
	}{
		{
			Name:   "no replacement",
			Branch: "feature-wip",
		},
		{
			Name:   "replace slashes",
			Branch: "feature/wip",
		},
		{
			Name:   "replace underscores",
			Branch: "feature_wip",
		},
	}

	for _, s := range tt {
		t.Run(s.Name, func(t *testing.T) {
			g := &Git{
				Branch: s.Branch,
			}
			got := g.SanitizedBranch()
			assert.Equal(t, want, got)
		})
	}
}
