package set

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
}

func (c *Controller) Set(_ context.Context, _ *logrus.Entry, param *Param) error {
	// find and read config
	files := param.Files
	if len(files) == 0 {
		// find templates
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
