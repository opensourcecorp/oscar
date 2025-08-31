package semver

import (
	"fmt"
	"regexp"
	"strings"

	iprint "github.com/opensourcecorp/oscar/internal/print"
	"golang.org/x/mod/semver"
)

// GetSemver tries to build a (Go) compliant Semantic Version number out of the provided string,
// regardless of how dirty it is. Despite using the semver package in a few  places internally, most
// of this implementation is custom due to limitations in that package -- like not being able to
// parse out just the Patch number, or Pre-Release/Build numbers not being allowed in
// semver.Canonical()
func GetSemver(s string) (string, error) {
	var v, preRelease, build string

	// Grab the semver parts separately so we can clean them up. Firstly, the Major-Minor-Patch
	// parts, but Patch takes some extra work to suss out -- if there's any prerelease or build
	// parts of the version, these show up in the "patch" index if we just split on dots, so we can
	// just grab the whole MMP with a regex
	v = "v" + regexp.MustCompile(`\d+(\.\d+)?(\.\d+)?`).FindStringSubmatch(s)[0]
	v = semver.Canonical(v)

	// Gross
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

	if !semver.IsValid(v) {
		return "", fmt.Errorf("could not understand the semantic version you provided in your Oscarfile: '%s'", s)
	}

	if v != s {
		iprint.Warnf("the Semantic Version string built was different from the one provided -- please edit your version to match the correct format: %s --> %s\n", s, v)
	}

	return v, nil
}
