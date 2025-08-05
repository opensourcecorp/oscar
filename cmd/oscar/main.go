// Package main runs oscar.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	icli "github.com/opensourcecorp/oscar/internal/cli"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

// TODO: keep init() but move all this logic elsewhere
func init() {
	var initHomedirsErrs error
	dirs := []string{
		consts.OscarHome,
		consts.OscarHomeBin,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			initHomedirsErrs = errors.Join(initHomedirsErrs, err)
		}
	}
	if initHomedirsErrs != nil {
		iprint.Errorf(
			"Internal error(s) when creating oscar directory '%s': %v\n",
			consts.OscarHome, initHomedirsErrs,
		)
		os.Exit(1)
	}

	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", consts.OscarHomeBin, os.Getenv("PATH"))); err != nil {
		iprint.Errorf(
			"Internal error(s) when setting oscar home bin path '%s': %v\n",
			consts.OscarHomeBin, initHomedirsErrs,
		)
		os.Exit(1)
	}
}

func main() {
	if err := icli.NewRootCmd().Run(context.Background(), os.Args); err != nil {
		iprint.Errorf("running: %v\n", err)
		os.Exit(1)
	}
}
