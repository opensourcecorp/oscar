package oscarcfg

import (
	"encoding/json"
	"fmt"
	"os"

	"buf.build/go/protovalidate"
	"github.com/opensourcecorp/oscar/internal/consts"
	oscarcfgpbv1 "github.com/opensourcecorp/oscar/internal/generated/opensourcecorp/oscar/config/v1"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/semver"
	"google.golang.org/protobuf/encoding/protojson"
)

// Get returns a populated [Config] based on the oscar config file location. If `path` is not
// provided, it will default to looking in the calling directory.
func Get(pathOverride ...string) (*oscarcfgpbv1.Config, error) {
	path := consts.DefaultOscarCfgFileName

	// Handle the override so we can test this function, and use it in other ways (like checking the
	// main branch's version data)
	if len(pathOverride) > 0 {
		path = pathOverride[0]
	}

	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading oscar config file: %w", err)
	}
	iprint.Debugf("data read from oscar config file:\n%s\n", string(yamlData))

	jsonSweepMap := make(map[string]any)
	if err := yaml.Unmarshal(yamlData, jsonSweepMap); err != nil {
		panic(err)
	}
	iprint.Debugf("YAML data unmarshalled to map: %+v\n", jsonSweepMap)

	jsonData, err := json.Marshal(jsonSweepMap)
	if err != nil {
		panic(err)
	}
	iprint.Debugf("map data as JSON string: %s\n", string(jsonData))

	var cfg = &oscarcfgpbv1.Config{}
	if err := protojson.Unmarshal(jsonData, cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling oscar config file '%s': %w", path, err)
	}
	iprint.Debugf("proto message: %+v\n", cfg)

	if err := protovalidate.Validate(cfg); err != nil {
		return nil, fmt.Errorf("validating oscar config file '%s': %w", path, err)
	}

	return cfg, nil
}

// VersionHasBeenIncremented reports whether the newVersion is greater than the oldVersion.
func VersionHasBeenIncremented(newVersion string, oldVersion string) bool {
	compValue := semver.Compare("v"+newVersion, "v"+oldVersion)
	iprint.Debugf("semver comparison value: %d\n", compValue)

	return compValue > 0
}
