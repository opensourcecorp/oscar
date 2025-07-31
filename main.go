package main

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar/internal/ci"
)

var banner = `
   ____________________
 /____________________/|
|  _   _   _  _   _  |/|
| | | |_  |  |_| |_| |/|
| |_|  _| |_ | | | \ |/|
|____________________|/
`

func main() {
	fmt.Println(banner)
	if err := ci.Run(); err != nil {
		// just exit, because all the errors were already logged
		os.Exit(1)
	}
}
