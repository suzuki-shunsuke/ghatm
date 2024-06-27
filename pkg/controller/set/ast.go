package set

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type Position struct {
	Line   int
	Column int
}

func parseWorkflowAST(content []byte, jobNames map[string]struct{}) ([]*Position, error) {
	file, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse a workflow file as YAML: %w", err)
	}
	list := []*Position{}
	for _, doc := range file.Docs {
		arr, err := parseDocAST(doc, jobNames)
		if err != nil {
			return nil, err
		}
		if len(arr) == 0 {
			continue
		}
		list = append(list, arr...)
	}
	return list, nil
}

func parseDocAST(doc *ast.DocumentNode, jobNames map[string]struct{}) ([]*Position, error) {
	body, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		return nil, nil
	}
	// jobs:
	//   jobName:
	//     timeout-minutes: 10
	//     steps:
	jobsNode := findJobsNode(body.Values)
	if jobsNode == nil {
		return nil, nil
	}
	return parseDocValue(jobsNode, jobNames)
}

func findJobsNode(values []*ast.MappingValueNode) *ast.MappingValueNode {
	for _, value := range values {
		key, ok := value.Key.(*ast.StringNode)
		if !ok {
			continue
		}
		if key.Value == "jobs" {
			return value
		}
	}
	return nil
}

func parseDocValue(value *ast.MappingValueNode, jobNames map[string]struct{}) ([]*Position, error) {
	jobs, ok := value.Value.(*ast.MappingNode)
	if !ok {
		return nil, errors.New("jobs must be a map")
	}
	arr := make([]*Position, 0, len(jobs.Values))
	for _, job := range jobs.Values {
		pos, err := parseJobAST(job, jobNames)
		if err != nil {
			return nil, err
		}
		if pos == nil {
			continue
		}
		arr = append(arr, pos)
	}
	return arr, nil
}

func parseJobAST(value *ast.MappingValueNode, jobNames map[string]struct{}) (*Position, error) {
	jobNameNode, ok := value.Key.(*ast.StringNode)
	if !ok {
		return nil, errors.New("job name must be a string")
	}
	jobName := jobNameNode.Value
	if _, ok := jobNames[jobName]; !ok {
		return nil, nil //nolint:nilnil
	}
	fields, ok := value.Value.(*ast.MappingNode)
	if !ok {
		return nil, errors.New("job value must be a *ast.MappingNode")
	}
	if len(fields.Values) == 0 {
		return nil, errors.New("job doesn't have any field")
	}
	firstValue := fields.Values[0]
	pos := firstValue.Key.GetToken().Position
	return &Position{
		Line:   pos.Line - 1,
		Column: pos.Column,
	}, nil
}
