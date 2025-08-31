package semver

import (
	"testing"
)

func TestGetSemver(t *testing.T) {
	t.Run("Convert regular semver to be conformant", func(t *testing.T) {
		s := "1.0.0.0"
		want := "v1.0.0"
		got, err := GetSemver(s)
		if err != nil {
			t.Errorf("unexpected error")
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("Convert semver with prerelease info to be conformant", func(t *testing.T) {
		s := "1.0-alpha"
		want := "v1.0.0-alpha"
		got, err := GetSemver(s)
		if err != nil {
			t.Errorf("unexpected error")
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("Convert semver with build info to be conformant", func(t *testing.T) {
		s := "1.0.2+abc"
		want := "v1.0.2+abc"
		got, err := GetSemver(s)
		if err != nil {
			t.Errorf("unexpected error")
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})

	t.Run("No conversion on a conformant semver", func(t *testing.T) {
		s := "v1.1.9-prebeta1+abc"
		want := s
		got, err := GetSemver(s)
		if err != nil {
			t.Errorf("unexpected error")
		}
		if want != got {
			t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
		}
	})
}
