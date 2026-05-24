package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v88/github"
	"golang.org/x/oauth2"
)

type Response = github.Response

func newGitHub(ctx context.Context) (*github.Client, error) {
	client, err := github.NewClient(github.WithHTTPClient(getHTTPClientForGitHub(ctx, getGitHubToken())))
	if err != nil {
		return nil, fmt.Errorf("create a GitHub client: %w", err)
	}
	return client, nil
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

func NewClient(ctx context.Context) (*Client, error) {
	gh, err := newGitHub(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		actions: gh.Actions,
	}, nil
}
