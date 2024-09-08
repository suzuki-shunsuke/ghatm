package edit

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Position struct {
	JobKey string
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
		return nil, errors.New("document body must be *ast.MappingNode")
	}
	// jobs:
	//   jobName:
	//     timeout-minutes: 10
	//     steps:
	jobsNode := findJobsNode(body.Values)
	if jobsNode == nil {
		return nil, errors.New("the field 'jobs' is required")
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

func getMappingValueNodes(value *ast.MappingValueNode) ([]*ast.MappingValueNode, error) {
	switch node := value.Value.(type) {
	case *ast.MappingNode:
		return node.Values, nil
	case *ast.MappingValueNode:
		return []*ast.MappingValueNode{node}, nil
	}
	return nil, errors.New("value must be either a *ast.MappingNode or a *ast.MappingValueNode")
}

func parseDocValue(value *ast.MappingValueNode, jobNames map[string]struct{}) ([]*Position, error) {
	values, err := getMappingValueNodes(value)
	if err != nil {
		return nil, err
	}
	arr := make([]*Position, 0, len(values))
	for _, job := range values {
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
	fields, err := getMappingValueNodes(value)
	if err != nil {
		return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	if len(fields) == 0 {
		return nil, logerr.WithFields(errors.New("job doesn't have any field"), logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	firstValue := fields[0]
	pos := firstValue.Key.GetToken().Position
	return &Position{
		JobKey: jobName,
		Line:   pos.Line - 1,
		Column: pos.Column,
	}, nil
}
