package containertools

import (
	"context"
	"regexp"
	"testing"

	igit "github.com/opensourcecorp/oscar/internal/git"
	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstructImageURI(t *testing.T) {
	cfg, err := oscarcfg.Get("../../../oscarcfg/test.oscar.yaml")
	require.NoError(t, err)

	git, err := igit.New(context.Background())
	require.NoError(t, err)

	var want *regexp.Regexp
	if git.Branch == "main" {
		want = regexp.MustCompile(`ghcr.io/opensourcecorp/oscar:` + cfg.GetVersion())
	} else {
		want = regexp.MustCompile(`ghcr.io/opensourcecorp/oscar:.*-.*-.*`)
	}

	got, err := constructImageURI(context.Background(), cfg)
	require.NoError(t, err)

	assert.Regexp(t, want, got)
}
