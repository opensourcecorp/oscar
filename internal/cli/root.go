package icli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/opensourcecorp/oscar/internal/consts"
	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/tasks/ci"
	"github.com/opensourcecorp/oscar/internal/tasks/delivery"
	"github.com/urfave/cli/v3"
)

const (
	// Command names and their flags
	rootCmdName   = "oscar"
	debugFlagName = "debug"

	ciCommandName = "ci"

	deliverCommandName = "deliver"
)

// NewRootCmd defines & returns the CLI command used as oscar's entrypoint.
func NewRootCmd() *cli.Command {
	version, err := getVersion()
	if err != nil {
		iprint.Errorf("determining your version: %v\n", err)
		os.Exit(1)
	}

	cmd := &cli.Command{
		Name:    rootCmdName,
		Usage:   "The OpenSourceCorp Automation Runner",
		Version: version,
		Action:  rootAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    debugFlagName,
				Usage:   "Whether to print debug logs during oscar runs",
				Sources: cli.EnvVars(consts.DebugEnvVarName),
			},
		},
		Commands: []*cli.Command{
			{
				Name:   ciCommandName,
				Usage:  "Runs CI tasks",
				Action: ciAction,
			},
			{
				Name:   deliverCommandName,
				Usage:  "Runs Delivery tasks",
				Action: deliverAction,
			},
		},
	}

	return cmd
}

// maybeSetDebug conditionally sets oscar's debug env var, so that other packages can use it.
func maybeSetDebug(cmd *cli.Command) {
	if cmd.Bool(debugFlagName) || os.Getenv(consts.DebugEnvVarName) != "" {
		_ = os.Setenv(consts.DebugEnvVarName, "true")
	}
}

// getVersion retrieves the version of the codebase.
func getVersion() (string, error) {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return "", fmt.Errorf("reading oscar config file: %w", err)
	}

	return cfg.Version, nil
}

// rootAction defines the logic for oscar's root command.
func rootAction(_ context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Debugf("oscar root command\n")
	_ = cli.ShowAppHelp(cmd)
	return errors.New("\nERROR: oscar requires a valid subcommand")
}

// ciAction defines the logic for oscar's ci subcommand.
func ciAction(ctx context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Banner()
	iprint.Debugf("oscar ci subcommand\n")

	if err := ci.Run(ctx); err != nil {
		return fmt.Errorf("running CI tasks: %w", err)
	}

	return nil
}

// deliverAction defines the logic for oscar's deliver subcommand.
func deliverAction(ctx context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Banner()
	iprint.Debugf("oscar deliver subcommand\n")

	if err := delivery.Run(ctx); err != nil {
		return fmt.Errorf("running Delivery tasks: %w", err)
	}

	return nil
}
