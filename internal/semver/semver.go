package semver

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tools"
	xsemver "golang.org/x/mod/semver"
)

// Get tries to build a compliant Semantic Version number out of the provided string, regardless of
// how dirty it is. Despite using the "golang.org/x/mod/semver" package in a few places internally,
// most of this implementation is custom due to limitations in that package -- like not being able
// to parse out just the Patch number, or Pre-Release/Build numbers not being allowed in
// semver.Canonical()
func Get(s string) (string, error) {
	// Grab the semver parts separately so we can clean them up. Firstly, the Major-Minor-Patch
	// parts, but Patch takes some extra work to suss out -- if there's any prerelease or build
	// parts of the version, these show up in the "patch" index if we just split on dots, so we can
	// just grab the whole MMP with a regex
	matchList := regexp.MustCompile(`[0-9]+(\.[0-9]+)?(\.[0-9]+)?`).FindStringSubmatch(s)
	if len(matchList) == 0 {
		return "", fmt.Errorf("malformed or unmatchable Semantic Version number from 'VERSION' file (got: '%s')", s)
	}
	v := matchList[0]

	// NOTE: the external semver package has some niceties, and we use them here, but since it's a
	// Go package it expects a "v" prefix on every number. We want to just keep the non-"v" data
	// since it's more portable, so we need to self-prefix the version number for the remaining
	// duration of this function.
	v = xsemver.Canonical("v" + v)
	if v == "" {
		return "", fmt.Errorf("unable to canonicalize provided version '%s' (after possibly converting to '%s')", s, v)
	}

	// Gross.
	var preRelease, build string
	prSplit := strings.Split(s, "-")
	bSplit := strings.Split(s, "+")
	if len(prSplit) > 1 {
		// could still have a build number
		preRelease = strings.Split(prSplit[1], "+")[0]
	}
	if len(bSplit) > 1 {
		build = bSplit[1]
	}

	if preRelease != "" {
		v += "-" + preRelease
	}
	if build != "" {
		v += "+" + build
	}

	if !xsemver.IsValid(v) {
		return "", fmt.Errorf("could not understand the Semantic Version you provided in your 'VERSION' file (got: '%s', converted to: '%s')", s, v)
	}

	// NOW, we can finally strip off the "v" prefix
	v = strings.TrimPrefix(v, "v")

	return v, nil
}

// GetFromFile reads the contents of the codebase's "VERSION" file, and passes it to [Get].
func GetFromFile(pathOverride ...string) (string, error) {
	root, err := tools.RunCommand([]string{"git", "rev-parse", "--show-toplevel"})
	if err != nil {
		return "", err
	}

	path := filepath.Join(root, "VERSION")
	if len(pathOverride) > 0 {
		path = pathOverride[0]
	}

	versionFileContents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	versionFileLines := strings.Split(string(versionFileContents), "\n")
	var version string
	for _, line := range versionFileLines {
		if v, err := Get(line); err == nil {
			version = v
			break
		}
		iprint.Debugf("error matching line in VERSION file: %v\n", err)
	}

	if version == "" {
		return "", fmt.Errorf("could not determine a Semantic Version from your 'VERSION' file")
	}

	return version, nil
}

// VersionWasIncremented reports whether the newVersion is greater than the oldVersion.
func VersionWasIncremented(newVersion string, oldVersion string) bool {
	compValue := xsemver.Compare("v"+newVersion, "v"+oldVersion)
	iprint.Debugf("semver comparison value: %d\n", compValue)

	return compValue > 0
}
