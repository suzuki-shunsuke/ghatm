package set

type Workflow struct {
	FilePath string `yaml:"-"`
	Jobs     map[string]*Job
}

type Job struct {
	Steps          []*Step
	Uses           string
	TimeoutMinutes int `yaml:"timeout-minutes"`
}

type Step struct {
	TimeoutMinutes int `yaml:"timeout-minutes"`
}
