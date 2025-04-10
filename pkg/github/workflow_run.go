package github

import (
	"context"

	"github.com/google/go-github/v71/github"
)

type ListWorkflowRunsOptions struct {
	Status string
	Page   int
}

type WorkflowRun struct {
	ID int64
	// Can be one of: completed, action_required, cancelled, failure, neutral, skipped, stale, success, timed_out, in_progress, queued, requested, waiting, pending
	Status string
}

func (c *Client) ListWorkflowRuns(ctx context.Context, owner, repo, workflowFileName string, opts *ListWorkflowRunsOptions) ([]*WorkflowRun, *github.Response, error) {
	o := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100, //nolint:mnd
			Page:    opts.Page,
		},
		Status: opts.Status,
	}
	runs, resp, err := c.actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFileName, o)
	if err != nil {
		return nil, resp, err //nolint:wrapcheck
	}
	ret := make([]*WorkflowRun, 0, len(runs.WorkflowRuns))
	for _, run := range runs.WorkflowRuns {
		ret = append(ret, &WorkflowRun{
			ID:     run.GetID(),
			Status: run.GetStatus(),
		})
	}
	return ret, resp, nil
}
