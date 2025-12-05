package set

import (
	"context"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/github"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
	Auto           bool
	RepoOwner      string
	RepoName       string
	Size           int
}

func Set(ctx context.Context, logger *slog.Logger, fs afero.Fs, param *Param) error {
	files := param.Files
	if len(files) == 0 {
		a, err := findWorkflows(fs)
		if err != nil {
			return err
		}
		files = a
	}

	var gh *github.Client
	if param.Auto {
		gh = github.NewClient(ctx)
	}

	for _, file := range files {
		logger := logger.With("workflow_file", file)
		logger.Info("handling the workflow file")
		if err := handleWorkflow(ctx, logger, fs, gh, file, param); err != nil {
			return slogerr.With(err, "workflow_file", file) //nolint:wrapcheck
		}
	}
	return nil
}
