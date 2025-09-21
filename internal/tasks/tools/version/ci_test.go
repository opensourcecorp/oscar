package versiontools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeGitRemote(t *testing.T) {
	remote := "git@github.com:opensourcecorp/oscar.git"
	want := "https://github.com/opensourcecorp/oscar.git"
	got := canonicalizeGitRemote(remote)
	assert.Equal(t, want, got)
}
