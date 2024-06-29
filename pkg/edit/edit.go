package edit

import (
	"bufio"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Edit(content []byte, timeout int) ([]byte, error) {
	wf := &Workflow{}
	if err := yaml.Unmarshal(content, wf); err != nil {
		return nil, fmt.Errorf("unmarshal a workflow file: %w", err)
	}
	if err := wf.Validate(); err != nil {
		return nil, fmt.Errorf("validate a workflow: %w", err)
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

func listJobsWithoutTimeout(jobs map[string]*Job) map[string]struct{} {
	m := make(map[string]struct{}, len(jobs))
	for jobName, job := range jobs {
		if hasTimeout(job) {
			continue
		}
		m[jobName] = struct{}{}
	}
	return m
}

func hasTimeout(job *Job) bool {
	if job.TimeoutMinutes != 0 || job.Uses != "" {
		return true
	}
	for _, step := range job.Steps {
		if step.TimeoutMinutes == 0 {
			return false
		}
	}
	return true
}

func insertTimeout(content []byte, positions []*Position, timeout int) ([]string, error) {
	reader := strings.NewReader(string(content))
	scanner := bufio.NewScanner(reader)
	num := -1

	lines := []string{}
	pos := positions[0]
	lastPosIndex := len(positions) - 1
	posIndex := 0
	for scanner.Scan() {
		num++
		line := scanner.Text()
		if pos.Line == num {
			indent := strings.Repeat(" ", pos.Column-1)
			lines = append(lines, indent+fmt.Sprintf("timeout-minutes: %d", timeout))
			if posIndex == lastPosIndex {
				pos.Line = -1
			} else {
				posIndex++
				pos = positions[posIndex]
			}
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan a workflow file: %w", err)
	}
	return lines, nil
}
