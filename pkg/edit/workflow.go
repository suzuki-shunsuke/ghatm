package edit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Workflow struct {
	Jobs map[string]*Job
}

type Job struct {
	Name           string
	Steps          []*Step
	Uses           string
	TimeoutMinutes any `yaml:"timeout-minutes"`
	Strategy       any
}

type Step struct {
	TimeoutMinutes any `yaml:"timeout-minutes"`
}

// foo (${{inputs.name}}) -> ^foo (.+?)$

var parameterRegexp = regexp.MustCompile(`\${{.+?}}`)

func (j *Job) GetName(k string) (string, *regexp.Regexp, error) {
	if j.Strategy == nil {
		if j.Name == "" {
			return k, nil, nil
		}
		if !strings.Contains(j.Name, "${{") {
			return j.Name, nil, nil
		}
		r, err := regexp.Compile("^" + parameterRegexp.ReplaceAllLiteralString(j.Name, ".+") + "$")
		if err != nil {
			return "", nil, fmt.Errorf("convert a job name with parameters to a regular expression: %w", err)
		}
		return j.Name, r, nil
	}
	if j.Name == "" {
		r, err := regexp.Compile("^" + k + ` \(.*\)$`)
		if err != nil {
			return "", nil, fmt.Errorf("convert a job name with matrix to a regular expression: %w", err)
		}
		return k, r, nil
	}
	if !strings.Contains(j.Name, "${{") {
		r, err := regexp.Compile("^" + j.Name + ` \(.*\)$`)
		if err != nil {
			return "", nil, fmt.Errorf("convert a job name with matrix to a regular expression: %w", err)
		}
		return j.Name, r, nil
	}
	r, err := regexp.Compile("^" + parameterRegexp.ReplaceAllLiteralString(j.Name, ".+") + "$")
	if err != nil {
		return "", nil, fmt.Errorf("convert a job name with parameters to a regular expression: %w", err)
	}
	return j.Name, r, nil
}

func (w *Workflow) Validate() error {
	if w == nil {
		return errors.New("workflow is nil")
	}
	if len(w.Jobs) == 0 {
		return errors.New("jobs are empty")
	}
	for jobName, job := range w.Jobs {
		if err := job.Validate(); err != nil {
			return logerr.WithFields(err, logrus.Fields{"job": jobName}) //nolint:wrapcheck
		}
	}
	return nil
}

func (j *Job) Validate() error {
	if j == nil {
		return errors.New("job is nil")
	}
	for _, step := range j.Steps {
		if err := step.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Step) Validate() error {
	if s == nil {
		return errors.New("step is nil")
	}
	return nil
}
