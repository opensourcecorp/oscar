// Package main runs oscar.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar"
	icli "github.com/opensourcecorp/oscar/internal/cli"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
)

func main() {
	if err := run(); err != nil {
		iprint.Errorf("running: %v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	if err := os.MkdirAll(consts.OscarHome, 0755); err != nil {
		return fmt.Errorf(
			"internal error when creating oscar home directory '%s': %v",
			consts.OscarHome, err,
		)
	}

	for name, value := range consts.MiseVars {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf(
				"internal error when setting mise env var '%s': %v",
				name, err,
			)
		}
	}

	cfgFileContents, err := oscar.Files.ReadFile("mise.toml")
	if err != nil {
		return fmt.Errorf("reading embedded file contents: %w", err)
	}

	if err := os.WriteFile(consts.MiseConfigFileName, cfgFileContents, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	defer func() {
		if rmErr := os.Remove(consts.MiseConfigFileName); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing mise config file: %w", rmErr))
		}
	}()

	if err := icli.NewRootCmd().Run(context.Background(), os.Args); err != nil {
		return fmt.Errorf("running: %v", err)
	}

	return err
}
