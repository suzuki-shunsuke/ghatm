package cli

import (
	"context"
	"io"
	"log/slog"

	"github.com/suzuki-shunsuke/urfave-cli-v3-util/helpall"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/vcmd"
	"github.com/urfave/cli/v3"
)

type Runner struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	LDFlags     *LDFlags
	Logger      *slog.Logger
	LogLevelVar *slog.LevelVar
}

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (r *Runner) Run(ctx context.Context, args ...string) error {
	return helpall.With(&cli.Command{ //nolint:wrapcheck
		Name:    "ghatm",
		Usage:   "",
		Version: r.LDFlags.Version + " (" + r.LDFlags.Commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level",
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			(&setCommand{
				logger:      r.Logger,
				logLevelVar: r.LogLevelVar,
			}).command(),
			(&completionCommand{
				stdout: r.Stdout,
			}).command(),
			vcmd.New(&vcmd.Command{
				Name:    "ghatm",
				Version: r.LDFlags.Version,
				SHA:     r.LDFlags.Commit,
			}),
		},
	}, nil).Run(ctx, args)
}
