package set

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func (c *Controller) handleWorkflow(file string, timeout int) error {
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}
	wf := &Workflow{}
	if err := yaml.Unmarshal(b, wf); err != nil {
		return fmt.Errorf("unmarshal a workflow file: %w", err)
	}
	jobNames := listJobsWithoutTimeout(wf.Jobs)
	positions, err := parseWorkflowAST(b, jobNames)
	if err != nil {
		return err
	}
	if len(positions) == 0 {
		return nil
	}
	fixLines(positions)

	lines, err := insertTimeout(b, positions, timeout)
	if err != nil {
		return err
	}

	return c.writeWorkflow(file, lines)
}

func fixLines(positions []*Position) {
	for i, pos := range positions {
		pos.Line += i
	}
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
				continue
			}
			posIndex++
			pos = positions[posIndex]
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan a workflow file: %w", err)
	}
	return lines, nil
}

func (c *Controller) writeWorkflow(file string, lines []string) error {
	stat, err := c.fs.Stat(file)
	if err != nil {
		return fmt.Errorf("get the workflow file stat: %w", err)
	}

	if err := afero.WriteFile(c.fs, file, []byte(strings.Join(lines, "\n")+"\n"), stat.Mode()); err != nil {
		return fmt.Errorf("write the workflow file: %w", err)
	}
	return nil
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
