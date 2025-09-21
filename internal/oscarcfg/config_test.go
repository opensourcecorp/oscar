package oscarcfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfigFilePath = "test.oscar.yaml"

func TestGet(t *testing.T) {
	cfg, err := Get(testConfigFilePath)
	require.NoError(t, err)

	t.Logf("parsed cfg:\n< %+v >", cfg)

	t.Run("version", func(t *testing.T) {
		want := "1.0.0"
		assert.Equal(t, want, cfg.GetVersion())
	})

	t.Run("deliver", func(t *testing.T) {
		wantBuildSources := []string{"./cmd/test"}
		gotBuildSources := cfg.GetDeliverables().GetGoGithubRelease().GetBuildSources()

		assert.Equal(t, wantBuildSources, gotBuildSources)
	})
}
