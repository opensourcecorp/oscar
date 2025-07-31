package iprint

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar/internal/consts"
)

// Errorf is a helper function that writes to standard error.
func Errorf(format string, args ...any) {
	if _, err := fmt.Fprintf(os.Stderr, format, args...); err != nil {
		// NOTE: panicking is fine here, this would be catastrophic lol
		panic(fmt.Sprintf("trying to write to stderr: %v", err))
	}
}

// Debugf is a helper function that prints debug logs if requested.
func Debugf(format string, args ...any) {
	if os.Getenv(consts.DebugEnvVarName) != "" {
		format = "DEBUG: " + format
		fmt.Printf(format, args...)
	}
}

// Banner prints the oscar banner.
func Banner() {
	var banner = `
   ____________________
 /____________________/|
|  _   _   _  _   _  |/|
| | | |_  |  |_| |_| |/|
| |_|  _| |_ | | | \ |/|
|____________________|/
`
	fmt.Println(banner)
}
