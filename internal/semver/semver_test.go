package semver

import (
	"testing"
)

func TestGetSemver(t *testing.T) {
	t.Run("No conversion on a basic conformant semver", func(t *testing.T) {
		s := "1.0.0"
		want := "1.0.0"
		got, err := Get(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("No conversion on a full conformant semver", func(t *testing.T) {
		s := "1.1.9-prebeta1+abc"
		want := s
		got, err := Get(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("Convert semver with prerelease info to be conformant", func(t *testing.T) {
		s := "1.0-alpha"
		want := "1.0.0-alpha"
		got, err := Get(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("Convert semver with build info to be conformant", func(t *testing.T) {
		s := "1.0+abc"
		want := "1.0.0+abc"
		got, err := Get(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})
}
