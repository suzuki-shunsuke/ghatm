package cli

import (
	"context"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:  "ghatm",
		Usage: "Set timeout-minutes to all GitHub Actions jobs. https://github.com/suzuki-shunsuke/ghatm",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level",
			},
			&cli.StringFlag{
				Name:  "log-color",
				Usage: "Log color. One of 'auto', 'always' (default), 'never'",
			},
		},
		Commands: []*cli.Command{
			(&setCommand{}).command(logger),
		},
	}).Run(ctx, env.Args)
}
