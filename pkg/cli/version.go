package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

type versionCommand struct{}

func (vc *versionCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show version",
		Action: vc.action,
	}
}

func (vc *versionCommand) action(_ context.Context, cmd *cli.Command) error {
	cli.ShowVersion(cmd)
	return nil
}
