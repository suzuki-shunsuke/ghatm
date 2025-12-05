package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/controller/set"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type setCommand struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
}

func (rc *setCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes",
		UsageText: "ghatm set",
		Description: `Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes.

$ ghatm set
`,
		Action: rc.action,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "timeout-minutes",
				Aliases: []string{"t"},
				Usage:   "The value of timeout-minutes",
				Value:   30, //nolint:mnd
			},
			&cli.BoolFlag{
				Name:    "auto",
				Aliases: []string{"a"},
				Usage:   "Estimate the value of timeout-minutes automatically",
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"r"},
				Usage:   "GitHub Repository",
				Sources: cli.EnvVars("GITHUB_REPOSITORY"),
			},
			&cli.IntFlag{
				Name:    "size",
				Aliases: []string{"s"},
				Usage:   "Data size",
				Value:   30, //nolint:mnd
			},
		},
	}
}

func (rc *setCommand) action(ctx context.Context, cmd *cli.Command) error {
	fs := afero.NewOsFs()
	logger := rc.logger
	if err := slogutil.SetLevel(rc.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	repo := cmd.String("repo")
	param := &set.Param{
		Files:          cmd.Args().Slice(),
		TimeoutMinutes: cmd.Int("timeout-minutes"),
		Auto:           cmd.Bool("auto"),
		Size:           cmd.Int("size"),
	}
	if param.Auto && repo == "" {
		return errors.New("the flag -auto requires the flag -repo")
	}
	if repo != "" {
		owner, repoName, ok := strings.Cut(repo, "/")
		if !ok {
			return fmt.Errorf("split the repository name: %s", repo)
		}
		param.RepoOwner = owner
		param.RepoName = repoName
	}
	return set.Set(ctx, logger, fs, param) //nolint:wrapcheck
}
