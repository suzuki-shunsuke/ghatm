package set

import (
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
}

func (c *Controller) Set(param *Param) error {
	files := param.Files
	if len(files) == 0 {
		a, err := FindWorkflows(c.fs)
		if err != nil {
			return err
		}
		files = a
	}

	for _, file := range files {
		if err := c.handleWorkflow(file, param.TimeoutMinutes); err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file": file,
			})
		}
	}
	return nil
}
