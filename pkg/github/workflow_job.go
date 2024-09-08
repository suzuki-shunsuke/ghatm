package github

import (
	"context"
	"time"

	"github.com/google/go-github/v64/github"
)

type ListWorkflowJobsOptions struct {
	Date   int
	Status string
	Page   int
}

type WorkflowJob struct {
	ID   int64
	Name string
	// Can be one of: completed, action_required, cancelled, failure, neutral, skipped, stale, success, timed_out, in_progress, queued, requested, waiting, pending
	Status   string
	Duration time.Duration
}

func (c *Client) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListWorkflowJobsOptions) ([]*WorkflowJob, *github.Response, error) {
	o := &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100, //nolint:mnd
			Page:    opts.Page,
		},
	}
	jobs, resp, err := c.actions.ListWorkflowJobs(ctx, owner, repo, runID, o)
	if err != nil {
		return nil, resp, err //nolint:wrapcheck
	}
	ret := make([]*WorkflowJob, 0, len(jobs.Jobs))
	for _, job := range jobs.Jobs {
		s := job.GetStartedAt()
		started := s.GetTime()
		if started == nil {
			continue
		}
		ret = append(ret, &WorkflowJob{
			ID:       job.GetID(),
			Name:     job.GetName(),
			Status:   job.GetStatus(),
			Duration: job.GetCompletedAt().Sub(*started),
		})
	}
	return ret, resp, nil
}
