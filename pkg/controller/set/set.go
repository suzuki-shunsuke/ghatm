package set

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Param struct {
	Files          []string
	TimeoutMinutes int
}

func Set(fs afero.Fs, param *Param) error {
	files := param.Files
	if len(files) == 0 {
		a, err := findWorkflows(fs)
		if err != nil {
			return err
		}
		files = a
	}

	for _, file := range files {
		if err := handleWorkflow(fs, file, param.TimeoutMinutes); err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file": file,
			})
		}
	}
	return nil
}
