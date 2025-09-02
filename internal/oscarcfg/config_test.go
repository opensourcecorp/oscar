package oscarcfg

import (
	"testing"
)

func TestRead(t *testing.T) {
	cfg, err := Get("test.oscar.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("parsed cfg: %+v", cfg)

	t.Run("version", func(t *testing.T) {
		want := "1.0.0"
		if cfg.Version != want {
			t.Errorf("version: wanted '%s', got '%s'", want, cfg.Version)
		}
	})

	t.Run("deliver", func(t *testing.T) {
		wantGHRelease := "test"
		gotGHRelease := cfg.Deliver.GoGitHubRelease.Repo
		if gotGHRelease != wantGHRelease {
			t.Errorf("GH Release: wanted '%s', got '%s'", wantGHRelease, gotGHRelease)
		}
	})
}
