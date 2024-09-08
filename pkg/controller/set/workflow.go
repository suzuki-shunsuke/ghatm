package set

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghatm/pkg/edit"
	"github.com/suzuki-shunsuke/ghatm/pkg/github"
	"gopkg.in/yaml.v3"
)

type GitHub interface {
	ListWorkflowRuns(ctx context.Context, owner, repo, workflowFileName string, opts *github.ListWorkflowRunsOptions) ([]*github.WorkflowRun, *github.Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *github.ListWorkflowJobsOptions) ([]*github.WorkflowJob, *github.Response, error)
}

func handleWorkflow(ctx context.Context, logE *logrus.Entry, fs afero.Fs, gh GitHub, file string, param *Param) error {
	content, err := afero.ReadFile(fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}

	wf := &edit.Workflow{}
	if err := yaml.Unmarshal(content, wf); err != nil {
		return fmt.Errorf("unmarshal a workflow file: %w", err)
	}
	if err := wf.Validate(); err != nil {
		return fmt.Errorf("validate a workflow: %w", err)
	}

	jobNames := edit.ListJobsWithoutTimeout(wf.Jobs)

	var timeouts map[string]int
	if param.Auto {
		tm, err := estimateTimeout(ctx, logE, gh, param, file, wf, jobNames)
		if err != nil {
			return err
		}
		timeouts = tm
	}

	after, err := edit.Edit(content, wf, timeouts, param.TimeoutMinutes)
	if err != nil {
		return fmt.Errorf("create a new workflow content: %w", err)
	}
	if after == nil {
		return nil
	}
	return writeWorkflow(fs, file, after)
}

func writeWorkflow(fs afero.Fs, file string, content []byte) error {
	stat, err := fs.Stat(file)
	if err != nil {
		return fmt.Errorf("get the workflow file stat: %w", err)
	}

	if err := afero.WriteFile(fs, file, content, stat.Mode()); err != nil {
		return fmt.Errorf("write the workflow file: %w", err)
	}
	return nil
}

func findWorkflows(fs afero.Fs) ([]string, error) {
	files, err := afero.Glob(fs, ".github/workflows/*.yml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yml: %w", err)
	}
	files2, err := afero.Glob(fs, ".github/workflows/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yaml: %w", err)
	}
	return append(files, files2...), nil
}
