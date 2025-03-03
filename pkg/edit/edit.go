package edit

import (
	"bufio"
	"fmt"
	"strings"
)

func Edit(content []byte, wf *Workflow, timeouts map[string]int, timeout int) ([]byte, error) {
	jobNames := ListJobsWithoutTimeout(wf.Jobs)
	positions, err := parseWorkflowAST(content, jobNames)
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 {
		return nil, nil
	}

	lines, err := insertTimeout(content, positions, timeouts, timeout)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}

func ListJobsWithoutTimeout(jobs map[string]*Job) map[string]struct{} {
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
	if job.Uses != "" {
		return true
	}

	switch v := job.TimeoutMinutes.(type) {
	case int:
		if v != 0 {
			return true
		}
	// An expression like "${{ inputs.timeout }}" is considered as having a timeout
	case string:
		if strings.Contains(v, "${{") {
			return true
		}
	case nil:
		// TimeoutMinutes is not set
	default:
		// Any other non-nil value is considered as having a timeout for future compatibility
		return true
	}

	for _, step := range job.Steps {
		if step.TimeoutMinutes == 0 {
			return false
		}
	}
	return true
}

func getTimeout(timeouts map[string]int, timeout int, jobKey string) int {
	if timeouts == nil {
		return timeout
	}
	if a, ok := timeouts[jobKey]; ok {
		return a
	}
	return -1
}

func insertTimeout(content []byte, positions []*Position, timeouts map[string]int, timeout int) ([]string, error) {
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
			if t := getTimeout(timeouts, timeout, pos.JobKey); t != -1 {
				lines = append(lines, indent+fmt.Sprintf("timeout-minutes: %d", t))
			}
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
