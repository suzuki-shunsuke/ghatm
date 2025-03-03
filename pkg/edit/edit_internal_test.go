package edit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEdit(t *testing.T) { //nolint:gocognit,cyclop,funlen
	t.Parallel()
	data := []struct {
		name     string
		content  string
		result   string
		isErr    bool
		wf       *Workflow
		timeouts map[string]int
	}{
		{
			name:    "normal",
			content: "normal.yaml",
			result:  "normal_result.yaml",
			wf: &Workflow{
				Jobs: map[string]*Job{
					"actionlint": {
						Uses: "suzuki-shunsuke/actionlint-workflow/.github/workflows/actionlint.yaml@813a6d08c08cfd7a08618a89a59bfe78e573597c # v1.0.1",
					},
					"foo": {
						TimeoutMinutes: 5,
						Steps: []*Step{
							{},
						},
					},
					"bar": {
						Steps: []*Step{
							{
								TimeoutMinutes: 5,
							},
						},
					},
					"zoo": {
						Steps: []*Step{
							{},
						},
					},
					"yoo": {
						Steps: []*Step{
							{},
						},
					},
				},
			},
		},
		{
			name:    "nochange",
			content: "nochange.yaml",
			wf: &Workflow{
				Jobs: map[string]*Job{
					"actionlint": {
						Uses: "suzuki-shunsuke/actionlint-workflow/.github/workflows/actionlint.yaml@813a6d08c08cfd7a08618a89a59bfe78e573597c # v1.0.1",
					},
					"foo": {
						TimeoutMinutes: 5,
						Steps: []*Step{
							{},
						},
					},
					"bar": {
						Steps: []*Step{
							{
								TimeoutMinutes: 5,
							},
						},
					},
				},
			},
		},
		{
			// The tool should recognize ${{ inputs.timeout }} in with-timeout job and add timeout to without-timeout job
			name:    "reusable_workflow_timeout",
			content: "reusable_workflow_timeout.yaml",
			result:  "reusable_workflow_timeout_result.yaml",
			wf: &Workflow{
				Jobs: map[string]*Job{
					"with-timeout": {
						TimeoutMinutes: "${{ inputs.timeout }}", // This should be detected as having a timeout via inputs
						Steps: []*Step{
							{},
						},
					},
					"without-timeout": {
						TimeoutMinutes: nil, // This should get a default timeout added
						Steps: []*Step{
							{},
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			content, err := os.ReadFile(filepath.Join("testdata", d.content))
			if err != nil {
				t.Fatal(err)
			}
			var expResult []byte
			if d.result != "" {
				content, err := os.ReadFile(filepath.Join("testdata", d.result))
				if err != nil {
					t.Fatal(err)
				}
				expResult = content
			}
			result, err := Edit(content, d.wf, d.timeouts, 30)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if result == nil {
				if expResult == nil {
					return
				}
				t.Fatalf("wanted %v, got nil", string(expResult))
			}
			if expResult == nil {
				t.Fatalf("wanted nil, got %v", string(result))
			}
			if diff := cmp.Diff(string(expResult), string(result)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
