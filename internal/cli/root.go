package icli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/opensourcecorp/oscar"
	"github.com/opensourcecorp/oscar/internal/ci"
	"github.com/opensourcecorp/oscar/internal/consts"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/urfave/cli/v3"
)

const (
	// Command names and their flags
	rootCmdName   = "oscar"
	debugFlagName = "debug"

	ciCommandName = "ci"
)

// NewRootCmd defines & returns the CLI command used as oscar's entrypoint.
func NewRootCmd() *cli.Command {
	cmd := &cli.Command{
		Name:    rootCmdName,
		Usage:   "The OpenSourceCorp Automation Runner",
		Version: getVersion(),
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
		},
	}

	return cmd
}

func maybeSetDebug(cmd *cli.Command) {
	if cmd.Bool(debugFlagName) || os.Getenv(consts.DebugEnvVarName) != "" {
		_ = os.Setenv(consts.DebugEnvVarName, "true")
	}
}

func getVersion() string {
	contents, err := oscar.Files.ReadFile("VERSION")
	if err != nil {
		panic(fmt.Sprintf("Internal error trying to read VERSION file: %v", err))
	}

	splits := strings.Split(string(contents), "\n")
	var version string
	for _, line := range splits {
		if !strings.HasPrefix(line, "#") {
			version = line
			break
		}
	}
	return version
}

func rootAction(_ context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Debugf("oscar root command\n")
	msg := "\nERROR: oscar requires a subcommand"
	_ = cli.ShowAppHelp(cmd)
	iprint.Errorf("%s\n", msg)
	return errors.New(msg)
}

func ciAction(_ context.Context, cmd *cli.Command) error {
	maybeSetDebug(cmd)
	iprint.Banner()
	iprint.Debugf("oscar ci subcommand\n")

	if err := ci.Run(); err != nil {
		return err
	}

	return nil
}
