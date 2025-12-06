package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/controller/set"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type setCommand struct{}

func (rc *setCommand) command(logger *slogutil.Logger) *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes",
		UsageText: "ghatm set",
		Description: `Set timeout-minutes to GitHub Actions jobs which don't have timeout-minutes.

$ ghatm set
`,
		Action: urfave.Action(rc.action, logger),
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

func (rc *setCommand) action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	fs := afero.NewOsFs()
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
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
	return set.Set(ctx, logger.Logger, fs, param) //nolint:wrapcheck
}
