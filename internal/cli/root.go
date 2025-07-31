package icli

import (
	"context"
	"errors"
	"os"

	"github.com/opensourcecorp/oscar/internal/ci"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/urfave/cli/v3"
)

const (
	debugFlagName = "debug"
)

func NewRootCmd() *cli.Command {
	cmd := &cli.Command{
		Name:   "oscar",
		Usage:  "The OpenSourceCorp Automation Runner",
		Action: rootAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    debugFlagName,
				Usage:   "Whether to print debug logs during oscar runs",
				Sources: cli.EnvVars(consts.DebugEnvVarName),
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "ci",
				Usage:  "Runs CI tasks",
				Action: ciAction,
			},
		},
	}

	return cmd
}

func rootAction(ctx context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Debugf("oscar root command\n")
	msg := "oscar requires a subcommand. Run 'oscar --help' for help."
	iprint.Errorf("%s\n", msg)
	return errors.New(msg)
}

func ciAction(_ context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Debugf("oscar ci subcommand\n")
	if err := ci.Run(); err != nil {
		return err
	}

	return nil
}

func maybeSetDebug(cmd *cli.Command) {
	if cmd.Bool(debugFlagName) || os.Getenv(consts.DebugEnvVarName) != "" {
		_ = os.Setenv(consts.DebugEnvVarName, "true")
	}
}
