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
		name    string
		content string
		result  string
		isErr   bool
	}{
		{
			name:    "normal",
			content: "normal.yaml",
			result:  "normal_result.yaml",
		},
		{
			name:    "nochange",
			content: "nochange.yaml",
		},
		{
			name:    "unmarshal error",
			content: "unmarshal_error.yaml",
			isErr:   true,
		},
		{
			name:    "jobs not found",
			content: "jobs_not_found.yaml",
			isErr:   true,
		},
		{
			name:    "jobs must be map",
			content: "invalid_jobs.yaml",
			isErr:   true,
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
			result, err := Edit(content, 30)
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
