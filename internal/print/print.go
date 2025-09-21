package iprint

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar/internal/consts"
)

// Banner prints the oscar banner.
func Banner() {
	var banner = `
           ____________________         \
=~=~=~=~=/____________________/|---------\
~=~=~=~=|  _   _   _  _   _  |/|----------\
=~=~=~=~| | | |_  |  |_| |_| |/|----------/
~=~=~=~=| |_|  _| |_ | | | \ |/|---------/
        |____________________|/         /
`
	fmt.Println(banner)
}

// Debugf is a helper function that prints debug logs if requested.
func Debugf(format string, args ...any) {
	if os.Getenv(consts.DebugEnvVarName) != "" {
		format = "DEBUG: " + format
		fmt.Printf(format, args...)
	}
}

// Warnf is a helper function that writes warnings to standard error.
func Warnf(format string, args ...any) {
	if _, err := fmt.Fprintf(os.Stderr, "WARN: "+format, args...); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(fmt.Sprintf("trying to write warning to stderr: %v", err))
	}
}

// Errorf is a helper function that writes errors to standard error.
func Errorf(format string, args ...any) {
	if _, err := fmt.Fprintf(os.Stderr, format, args...); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(fmt.Sprintf("trying to write error to stderr: %v", err))
	}
}
