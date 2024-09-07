package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/controller/set"
	"github.com/suzuki-shunsuke/ghatm/pkg/log"
	"github.com/urfave/cli/v2"
)

type setCommand struct {
	logE *logrus.Entry
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
				EnvVars: []string{"GITHUB_REPOSITORY"},
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

func (rc *setCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	logE := rc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	repo := c.String("repo")
	param := &set.Param{
		Files:          c.Args().Slice(),
		TimeoutMinutes: c.Int("timeout-minutes"),
		Auto:           c.Bool("auto"),
		Size:           c.Int("size"),
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
	return set.Set(c.Context, fs, param) //nolint:wrapcheck
}
