package set

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Workflow struct {
	Jobs map[string]*Job
}

type Job struct {
	Steps          []*Step
	Uses           string
	TimeoutMinutes int `yaml:"timeout-minutes"`
}

type Step struct {
	TimeoutMinutes int `yaml:"timeout-minutes"`
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
