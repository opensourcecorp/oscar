// Package main runs oscar.
package main

import (
	"context"
	"os"

	icli "github.com/opensourcecorp/oscar/internal/cli"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

func main() {
	if err := icli.NewRootCmd().Run(context.Background(), os.Args); err != nil {
		iprint.Errorf("running: %v\n", err)
		os.Exit(1)
	}
}
