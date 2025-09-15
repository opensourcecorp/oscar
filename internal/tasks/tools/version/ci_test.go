package versiontools

import "testing"

func TestCanonicalizeGitRemote(t *testing.T) {
	remote := "git@github.com:opensourcecorp/oscar.git"
	want := "https://github.com/opensourcecorp/oscar.git"
	got := canonicalizeGitRemote(remote)

	if want != got {
		t.Errorf("\nwant: %v\ngot:  %v", want, got)
	}
}
