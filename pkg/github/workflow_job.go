package github

import (
	"context"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/sirupsen/logrus"
)

type ListWorkflowJobsOptions struct {
	Date   int
	Status string
	Page   int
}

type WorkflowJob struct {
	ID   int64
	Name string
	// The phase of the lifecycle that the job is currently in.
	// "queued", "in_progress", "completed", "waiting", "requested", "pending"
	Status string
	// The outcome of the job.
	// "success", "failure", "neutral", "cancelled", "skipped", "timed_out", "action_required",
	Conclusion string
	Duration   time.Duration
}

func (c *Client) ListWorkflowJobs(ctx context.Context, logE *logrus.Entry, owner, repo string, runID int64, opts *ListWorkflowJobsOptions) ([]*WorkflowJob, *github.Response, error) {
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
		j := &WorkflowJob{
			ID:         job.GetID(),
			Name:       job.GetName(),
			Status:     job.GetStatus(),
			Conclusion: job.GetConclusion(),
			Duration:   job.GetCompletedAt().Sub(*started),
		}
		if j.Status != "completed" || j.Conclusion != "success" {
			logE.WithFields(logrus.Fields{
				"status":     j.Status,
				"conclusion": j.Conclusion,
			}).Debug("skip the job")
			continue
		}
		ret = append(ret, j)
	}
	return ret, resp, nil
}
