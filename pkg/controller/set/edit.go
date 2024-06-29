package set

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func edit(content []byte, timeout int) ([]byte, error) {
	wf := &Workflow{}
	if err := yaml.Unmarshal(content, wf); err != nil {
		return nil, fmt.Errorf("unmarshal a workflow file: %w", err)
	}
	jobNames := listJobsWithoutTimeout(wf.Jobs)
	positions, err := parseWorkflowAST(content, jobNames)
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 {
		return nil, nil
	}

	lines, err := insertTimeout(content, positions, timeout)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}
