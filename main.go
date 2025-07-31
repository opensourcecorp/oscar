package main

import (
	"context"
	"os"

	icli "github.com/opensourcecorp/oscar/internal/cli"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

func main() {
	iprint.Banner()
	if err := icli.NewRootCmd().Run(context.Background(), os.Args); err != nil {
		// just exit, because all the errors were already logged
		os.Exit(1)
	}
}
