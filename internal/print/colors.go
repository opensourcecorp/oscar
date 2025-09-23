package iprint

// NOTE: there's a lot of extra-feeling code in here, but that's because of how these vars/functions
// end up being used. Since oscar checks for the '--no-color' flag at execution startup (and then
// sets the env var for it), but that env var is checked *at overall program start time* (i.e.
// before any flags are parsed/handled), we can't just set top-level global vars per color that
// check the env var -- we have to only return color codes on-demand, so the var is available to be
// checked at any point after the env var is assigned.
//
// All this is fine, it just might be goofy-looking.

import (
	"os"

	"github.com/opensourcecorp/oscar/internal/consts"
	"golang.org/x/term"
)

var (
	// ANSI color codes
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	white   = "\033[97m"
)

// AllColors holds All the possible ANSI color codes that can be used. See the topmost comment in
// this file for "why".
type AllColors struct {
	Reset      string
	Red        string
	Green      string
	Yellow     string
	Blue       string
	Magenta    string
	Cyan       string
	Gray       string
	White      string
	DebugColor string
	InfoColor  string
	WarnColor  string
	ErrorColor string
	GoodColor  string
}

// Colors returns a populated [AllColors]. It can be used to grab a conditional ANSI color code in
// message-printing.
//
// For example, to override a print function's default text color with, say, red, you could do
// something like this:
//
//	fmt.Println(iprint.Colors().Red + "oops" + iprint.Colors().Reset)
func Colors() AllColors {
	out := AllColors{
		Reset:   color(reset),
		Red:     color(red),
		Green:   color(green),
		Yellow:  color(yellow),
		Blue:    color(blue),
		Magenta: color(magenta),
		Cyan:    color(cyan),
		Gray:    color(gray),
		White:   color(white),

		DebugColor: color(blue),
		InfoColor:  color(cyan),
		WarnColor:  color(yellow),
		ErrorColor: color(red),
		GoodColor:  color(green),
	}

	return out
}

// color is passed an ANSI color code, and conditionally returns either it or an empty string in the
// case of colors being disabled.
func color(ansiCode string) string {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return ""
	}
	if noColor := os.Getenv(consts.OscarEnvVarNoColor); noColor != "" {
		return ""
	}

	return ansiCode
}
