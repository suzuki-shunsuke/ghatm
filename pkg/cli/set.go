package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/controller/set"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type SetFlags struct {
	*GlobalFlags

	TimeoutMinutes int
	Auto           bool
	Repo           string
	Size           int
	Args           []string
}

type setCommand struct{}

func (rc *setCommand) command(logger *slogutil.Logger, globalFlags *GlobalFlags) *cli.Command {
	flags := &SetFlags{GlobalFlags: globalFlags}
	return &cli.Command{
		Name:      "set",
		Usage:     "Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes",
		UsageText: "ghatm set",
		Description: `Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes.

$ ghatm set
`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flags.Args = cmd.Args().Slice()
			return rc.action(ctx, logger, flags)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level",
				Destination: &flags.LogLevel,
			},
			&cli.StringFlag{
				Name:        "log-color",
				Usage:       "Log color. One of 'auto', 'always' (default), 'never'",
				Destination: &flags.LogColor,
			},
			&cli.IntFlag{
				Name:        "timeout-minutes",
				Aliases:     []string{"t"},
				Usage:       "The value of timeout-minutes",
				Value:       30, //nolint:mnd
				Destination: &flags.TimeoutMinutes,
			},
			&cli.BoolFlag{
				Name:        "auto",
				Aliases:     []string{"a"},
				Usage:       "Estimate the value of timeout-minutes automatically",
				Destination: &flags.Auto,
			},
			&cli.StringFlag{
				Name:        "repo",
				Aliases:     []string{"r"},
				Usage:       "GitHub Repository",
				Sources:     cli.EnvVars("GITHUB_REPOSITORY"),
				Destination: &flags.Repo,
			},
			&cli.IntFlag{
				Name:        "size",
				Aliases:     []string{"s"},
				Usage:       "Data size",
				Value:       30, //nolint:mnd
				Destination: &flags.Size,
			},
		},
	}
}

func (rc *setCommand) action(ctx context.Context, logger *slogutil.Logger, flags *SetFlags) error {
	fs := afero.NewOsFs()
	if err := logger.SetLevel(flags.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(flags.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	param := &set.Param{
		Files:          flags.Args,
		TimeoutMinutes: flags.TimeoutMinutes,
		Auto:           flags.Auto,
		Size:           flags.Size,
	}
	if param.Auto && flags.Repo == "" {
		return errors.New("the flag -auto requires the flag -repo")
	}
	if flags.Repo != "" {
		owner, repoName, ok := strings.Cut(flags.Repo, "/")
		if !ok {
			return fmt.Errorf("split the repository name: %s", flags.Repo)
		}
		param.RepoOwner = owner
		param.RepoName = repoName
	}
	return set.Set(ctx, logger.Logger, fs, param) //nolint:wrapcheck
}
