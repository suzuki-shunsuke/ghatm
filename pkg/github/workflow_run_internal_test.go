package github

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v89/github"
)

type mockActionsService struct {
	runs *github.WorkflowRuns
	resp *github.Response
	err  error
}

func (m *mockActionsService) ListWorkflowRunsByFileName(_ context.Context, _, _, _ string, _ *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	return m.runs, m.resp, m.err
}

func (m *mockActionsService) ListWorkflowJobs(_ context.Context, _, _ string, _ int64, _ *github.ListWorkflowJobsOptions) (*github.Jobs, *github.Response, error) {
	return nil, nil, nil
}

func newResponse(statusCode int) *github.Response {
	return &github.Response{
		Response: &http.Response{StatusCode: statusCode},
	}
}

func TestClient_ListWorkflowRuns(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		actions ActionsService
		isErr   bool
		wantLen int
	}{
		{
			// A workflow that has never run has no workflow runs, so the GitHub
			// API returns 404 Not Found. This must be treated as "no runs"
			// instead of a fatal error.
			name: "workflow which has never run returns 404",
			actions: &mockActionsService{
				resp: newResponse(http.StatusNotFound),
				err: &github.ErrorResponse{
					Response: newResponse(http.StatusNotFound).Response,
					Message:  "Not Found",
				},
			},
			wantLen: 0,
		},
		{
			name: "workflow runs are returned",
			actions: &mockActionsService{
				runs: &github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{
						{},
						{},
					},
				},
				resp: newResponse(http.StatusOK),
			},
			wantLen: 2,
		},
		{
			name: "non-404 errors are propagated",
			actions: &mockActionsService{
				resp: newResponse(http.StatusInternalServerError),
				err:  errors.New("internal server error"),
			},
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{actions: d.actions}
			runs, _, err := c.ListWorkflowRuns(context.Background(), "owner", "repo", "workflow.yaml", &ListWorkflowRunsOptions{})
			if d.isErr {
				if err == nil {
					t.Fatal("an error must be returned")
				}
				return
			}
			if err != nil {
				t.Fatalf("an error must not be returned: %v", err)
			}
			if len(runs) != d.wantLen {
				t.Fatalf("the number of workflow runs must be %d but got %d", d.wantLen, len(runs))
			}
		})
	}
}
