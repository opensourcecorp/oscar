package iprint

import (
	"fmt"
	"os"
	"time"

	"github.com/opensourcecorp/oscar/internal/consts"
)

// Banner prints the oscar stylistic banner.
func Banner() {
	var banner = `
           ____________________         \
=~=~=~=~=/____________________/|---------\
~=~=~=~=|  _   _   _  _   _  |/|----------\
=~=~=~=~| | | |_  |  |_| |_| |/|----------/
~=~=~=~=| |_|  _| |_ | | | \ |/|---------/
        |____________________|/         /
`

	if os.Getenv(consts.OscarEnvVarNoBanner) == "" {
		Goodf(banner + "\n")
	}
}

// Debugf is a helper function that prints debug logs if requested.
func Debugf(format string, args ...any) {
	colors := Colors()
	if os.Getenv(consts.OscarEnvVarDebug) != "" {
		fmt.Printf(colors.DebugColor+"DEBUG: "+format+colors.Reset, args...)
	}
}

// Infof is a helper function that writes info-level text.
func Infof(format string, args ...any) {
	colors := Colors()
	fmt.Printf(colors.InfoColor+format+colors.Reset, args...)
}

// Warnf is a helper function that writes warnings to standard error.
func Warnf(format string, args ...any) {
	colors := Colors()
	if _, err := fmt.Fprintf(
		os.Stderr,
		colors.WarnColor+"WARN: "+format+colors.Reset,
		args...,
	); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(
			fmt.Sprintf(colors.ErrorColor+"trying to write warning to stderr: %v"+colors.Reset, err),
		)
	}
}

// Errorf is a helper function that writes errors to standard error.
func Errorf(format string, args ...any) {
	colors := Colors()
	if _, err := fmt.Fprintf(
		os.Stderr,
		colors.ErrorColor+format+colors.Reset,
		args...,
	); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(
			fmt.Sprintf(colors.ErrorColor+"trying to write error to stderr: %v"+colors.Reset, err),
		)
	}
}

// Goodf is a helper function that prints green info text indicating something went well
func Goodf(format string, args ...any) {
	colors := Colors()
	fmt.Printf(colors.GoodColor+format+colors.Reset, args...)
}

// RunDurationString returns a calculated duration used to indicate how long a particular Task (or
// set of Tasks) took to run.
func RunDurationString(t time.Time) string {
	return fmt.Sprintf("t: %s", time.Since(t).Round(time.Second/1000).String())
}
