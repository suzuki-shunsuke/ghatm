package github

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

type Response = github.Response

func newGitHub(ctx context.Context) *github.Client {
	return github.NewClient(getHTTPClientForGitHub(ctx, getGitHubToken()))
}

func getGitHubToken() string {
	if token := os.Getenv("GHATM_GITHUB_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}

func getHTTPClientForGitHub(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}

type ActionsService interface {
	ListWorkflowRunsByFileName(ctx context.Context, owner, repo, workflowFileName string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *github.ListWorkflowJobsOptions) (*github.Jobs, *github.Response, error)
}

type Client struct {
	actions ActionsService
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		actions: newGitHub(ctx).Actions,
	}
}
