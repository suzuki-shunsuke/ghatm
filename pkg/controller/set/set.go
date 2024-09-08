package set

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/github"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
	Auto           bool
	RepoOwner      string
	RepoName       string
	Size           int
}

func Set(ctx context.Context, logE *logrus.Entry, fs afero.Fs, param *Param) error {
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
		if err := handleWorkflow(ctx, logE, fs, gh, file, param); err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file": file,
			})
		}
	}
	return nil
}
