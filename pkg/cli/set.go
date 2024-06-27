package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/gha-set-timeout-minutes/pkg/controller/set"
	"github.com/suzuki-shunsuke/gha-set-timeout-minutes/pkg/log"
	"github.com/urfave/cli/v2"
)

type setCommand struct {
	logE *logrus.Entry
}

func (rc *setCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes",
		UsageText: "gha-set-timeout-minutes set",
		Description: `Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes.

$ gha-set-timeout-minutes set
`,
		Action: rc.action,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "timeout-minutes",
				Aliases: []string{"m"},
				Usage:   "The value of timeout-minutes",
				Value:   30, //nolint:mnd
			},
		},
	}
}

func (rc *setCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	ctrl := set.NewController(fs)
	logE := rc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	return ctrl.Set(&set.Param{ //nolint:wrapcheck
		Files:          c.Args().Slice(),
		TimeoutMinutes: c.Int("timeout-minutes"),
	})
}
