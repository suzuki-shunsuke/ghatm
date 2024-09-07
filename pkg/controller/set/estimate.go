package set

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghatm/pkg/edit"
	"github.com/suzuki-shunsuke/ghatm/pkg/github"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"golang.org/x/exp/slices"
)

func setNamePatterns(jobs map[string]*edit.Job, jobKeys map[string]struct{}, staticNames map[string]struct{}, namePatterns map[string]*regexp.Regexp) error {
	for jobKey, job := range jobs {
		if _, ok := jobKeys[jobKey]; !ok {
			continue
		}
		name, nameRegexp, err := job.GetName(jobKey)
		if err != nil {
			return fmt.Errorf("get a job name: %w", logerr.WithFields(err, logrus.Fields{
				"job_key": jobKey,
			}))
		}
		if nameRegexp == nil {
			staticNames[name] = struct{}{}
			continue
		}
		namePatterns[name] = nameRegexp
	}
	return nil
}

func handleJob(jobDurationMap map[string][]time.Duration, staticNames map[string]struct{}, namePatterns map[string]*regexp.Regexp, job *github.WorkflowJob) {
	if _, ok := staticNames[job.Name]; ok {
		a, ok := jobDurationMap[job.Name]
		if !ok {
			a = []time.Duration{}
		}
		a = append(a, job.Duration)
		jobDurationMap[job.Name] = a
		return
	}
	for name, nameRegexp := range namePatterns {
		if !nameRegexp.MatchString(job.Name) {
			continue
		}
		a, ok := jobDurationMap[name]
		if !ok {
			a = []time.Duration{}
		}
		a = append(a, job.Duration)
		jobDurationMap[name] = a
		return
	}
}

func handleWorkflowRun(ctx context.Context, gh GitHub, param *Param, jobDurationMap map[string][]time.Duration, staticNames map[string]struct{}, namePatterns map[string]*regexp.Regexp, runID int64) (bool, error) {
	jobOpts := &github.ListWorkflowJobsOptions{
		Status: "success",
	}
	for range 10 {
		if isCompleted(jobDurationMap, param.Size) {
			return true, nil
		}
		jobs, resp, err := gh.ListWorkflowJobs(ctx, param.RepoOwner, param.RepoName, runID, jobOpts)
		if err != nil {
			return false, fmt.Errorf("list workflow jobs: %w", logerr.WithFields(err, logrus.Fields{
				"workflow_run_id": runID,
			}))
		}
		for _, job := range jobs {
			if isCompleted(jobDurationMap, param.Size) {
				return true, nil
			}
			handleJob(jobDurationMap, staticNames, namePatterns, job)
		}
		if resp.NextPage == 0 {
			break
		}
		jobOpts.Page = resp.NextPage
	}
	return true, nil
}

func getJobsByAPI(ctx context.Context, gh GitHub, param *Param, file string, wf *edit.Workflow, jobKeys map[string]struct{}) (map[string][]time.Duration, error) {
	staticNames := make(map[string]struct{}, len(wf.Jobs))
	namePatterns := make(map[string]*regexp.Regexp, len(wf.Jobs))
	if err := setNamePatterns(wf.Jobs, jobKeys, staticNames, namePatterns); err != nil {
		return nil, err
	}

	jobDurationMap := make(map[string][]time.Duration, len(wf.Jobs))

	runOpts := &github.ListWorkflowRunsOptions{
		Status: "success",
	}
	for range 10 {
		runs, resp, err := gh.ListWorkflowRuns(ctx, param.RepoOwner, param.RepoName, file, runOpts)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %w", err)
		}
		for _, run := range runs {
			completed, err := handleWorkflowRun(ctx, gh, param, jobDurationMap, staticNames, namePatterns, run.ID)
			if err != nil {
				return nil, err
			}
			if completed {
				return jobDurationMap, nil
			}
		}
		if resp.NextPage == 0 {
			return jobDurationMap, nil
		}
		runOpts.Page = resp.NextPage
	}
	return jobDurationMap, nil
}

func isCompleted(jobDurationMap map[string][]time.Duration, size int) bool {
	for _, durations := range jobDurationMap {
		if len(durations) < size {
			return false
		}
	}
	return true
}

func estimateTimeout(ctx context.Context, gh GitHub, param *Param, file string, wf *edit.Workflow, jobKeys map[string]struct{}) (map[string]int, error) {
	fileName := filepath.Base(file)
	jobs, err := getJobsByAPI(ctx, gh, param, fileName, wf, jobKeys)
	if err != nil {
		return nil, err
	}
	m := make(map[string]int, len(jobs))
	for jobKey, durations := range jobs {
		maxDuration := slices.Max(durations)
		m[jobKey] = int(math.Ceil(maxDuration.Minutes())) + 10 //nolint:mnd
	}

	return m, nil
}
