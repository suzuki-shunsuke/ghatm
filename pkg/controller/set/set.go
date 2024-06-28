package set

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
}

func (c *Controller) Set(logE *logrus.Entry, param *Param) error {
	files := param.Files
	if len(files) == 0 {
		a, err := FindWorkflows(c.fs)
		if err != nil {
			return err
		}
		files = a
	}

	var once sync.Once
	var failed bool
	onceBody := func() {
		failed = true
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	var semaphore chan struct{}
	if len(files) > 10 { //nolint:mnd
		semaphore = make(chan struct{}, 10) //nolint:mnd
	}
	for _, file := range files {
		go func(file string) {
			defer wg.Done()
			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
				}()
			}
			if err := c.handleWorkflow(file, param.TimeoutMinutes); err != nil {
				logerr.WithError(logE, err).WithField("file", file).Error("handle a workflow")
				once.Do(onceBody)
			}
		}(file)
	}
	wg.Wait()
	if failed {
		return errors.New("failed to handle workflows")
	}
	return nil
}
