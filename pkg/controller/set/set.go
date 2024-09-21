package set

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/edit"
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

	workflowCalls := map[string]map[string][]int{}

	for _, file := range files {
		if err := handleWorkflow(ctx, logE, fs, gh, file, param, workflowCalls); err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file": file,
			})
		}
	}

	for workflowFileName, jobs := range workflowCalls {
		for jobKey, timeouts := range jobs {
			after, err := edit.Edit(content, wf, timeouts, param.TimeoutMinutes)
			if err != nil {
				return fmt.Errorf("create a new workflow content: %w", err)
			}
			if after == nil {
				return nil
			}
			return writeWorkflow(fs, file, after)
		}
	}

	return nil
}
