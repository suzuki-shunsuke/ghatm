package set

import (
	"context"
	"fmt"
	"maps"
	"math"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghatm/pkg/edit"
	"github.com/suzuki-shunsuke/ghatm/pkg/github"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func setNamePatterns(jobs map[string]*edit.Job, jobKeys map[string]struct{}, staticNames map[string]string, namePatterns map[string]*regexp.Regexp) error {
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
			staticNames[name] = jobKey
			continue
		}
		namePatterns[jobKey] = nameRegexp
	}
	return nil
}

func handleJob(logE *logrus.Entry, jobDurationMap map[string][]time.Duration, staticNames map[string]string, namePatterns map[string]*regexp.Regexp, job *github.WorkflowJob) {
	if jobKey, ok := staticNames[job.Name]; ok {
		logE.WithFields(logrus.Fields{
			"job_name": job.Name,
			"job_key":  jobKey,
		}).Debug("adding the job duration")
		a, ok := jobDurationMap[jobKey]
		if !ok {
			a = []time.Duration{}
		}
		a = append(a, job.Duration)
		jobDurationMap[jobKey] = a
		return
	}
	for jobKey, nameRegexp := range namePatterns {
		if !nameRegexp.MatchString(job.Name) {
			continue
		}
		logE.WithFields(logrus.Fields{
			"job_name":         job.Name,
			"job_key":          jobKey,
			"job_name_pattern": nameRegexp.String(),
		}).Debug("adding the job duration")
		a, ok := jobDurationMap[jobKey]
		if !ok {
			a = []time.Duration{}
		}
		a = append(a, job.Duration)
		jobDurationMap[jobKey] = a
		return
	}
	logE.WithFields(logrus.Fields{
		"job_name": job.Name,
	}).Debug("the job name doesn't match")
}

func handleWorkflowRun(ctx context.Context, logE *logrus.Entry, gh GitHub, param *Param, jobDurationMap map[string][]time.Duration, staticNames map[string]string, namePatterns map[string]*regexp.Regexp, runID int64) (bool, error) {
	jobOpts := &github.ListWorkflowJobsOptions{
		Status: "success",
	}
	for range 10 {
		if isCompleted(logE, jobDurationMap, param.Size) {
			return true, nil
		}
		jobs, resp, err := gh.ListWorkflowJobs(ctx, logE, param.RepoOwner, param.RepoName, runID, jobOpts)
		if err != nil {
			return false, fmt.Errorf("list workflow jobs: %w", logerr.WithFields(err, logrus.Fields{
				"workflow_run_id": runID,
			}))
		}
		logE.WithField("num_of_jobs", len(jobs)).Debug("list workflow jobs")
		for _, job := range jobs {
			logE := logE.WithFields(logrus.Fields{
				"job_name":     job.Name,
				"job_status":   job.Status,
				"job_duration": job.Duration,
			})
			if isCompleted(logE, jobDurationMap, param.Size) {
				logE.Debug("job has been completed")
				return true, nil
			}
			logE.Debug("handling the job")
			handleJob(logE, jobDurationMap, staticNames, namePatterns, job)
		}
		if resp.NextPage == 0 {
			break
		}
		jobOpts.Page = resp.NextPage
	}
	return false, nil
}

// getJobsByAPI gets each job's durations by the GitHub API.
// It returns a map of job key and durations.
func getJobsByAPI(ctx context.Context, logE *logrus.Entry, gh GitHub, param *Param, file string, wf *edit.Workflow, jobKeys map[string]struct{}) (map[string][]time.Duration, error) {
	// jobName -> jobKey
	staticNames := make(map[string]string, len(wf.Jobs))
	// jobKey -> regular expression of job name
	namePatterns := make(map[string]*regexp.Regexp, len(wf.Jobs))
	if err := setNamePatterns(wf.Jobs, jobKeys, staticNames, namePatterns); err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"static_names":  strings.Join(slices.Collect(maps.Keys(staticNames)), ", "),
		"name_patterns": strings.Join(slices.Collect(maps.Keys(namePatterns)), ", "),
	}).Debug("static names and name patterns")

	jobDurationMap := make(map[string][]time.Duration, len(wf.Jobs))
	for jobKey := range jobKeys {
		jobDurationMap[jobKey] = []time.Duration{}
	}

	runOpts := &github.ListWorkflowRunsOptions{
		Status: "success",
	}
	loopSize := int(math.Ceil(float64(param.Size) * 3.0 / 100)) //nolint:mnd
	for range loopSize {
		runs, resp, err := gh.ListWorkflowRuns(ctx, param.RepoOwner, param.RepoName, file, runOpts)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %w", err)
		}
		logE.WithField("num_of_runs", len(runs)).Debug("list workflow runs")
		for _, run := range runs {
			completed, err := handleWorkflowRun(ctx, logE, gh, param, jobDurationMap, staticNames, namePatterns, run.ID)
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

func isCompleted(logE *logrus.Entry, jobDurationMap map[string][]time.Duration, size int) bool {
	for jobKey, durations := range jobDurationMap {
		if len(durations) < size {
			logE.WithFields(logrus.Fields{
				"job_key":          jobKey,
				"param_size":       size,
				"num_of_durations": len(durations),
			}).Debug("the job hasn't been completed")
			return false
		}
	}
	return true
}

// estimateTimeout estimates each job's timeout-minutes.
// It returns a map of job key and timeout-minutes.
// If there is no job's duration, the job is excluded from the return value.
func estimateTimeout(ctx context.Context, logE *logrus.Entry, gh GitHub, param *Param, file string, wf *edit.Workflow, jobKeys map[string]struct{}) (map[string]int, error) {
	fileName := filepath.Base(file)
	jobs, err := getJobsByAPI(ctx, logE, gh, param, fileName, wf, jobKeys)
	if err != nil {
		return nil, err
	}

	// Each job's timeout-minutes is `max(durations) + 10`.
	m := make(map[string]int, len(jobs))
	for jobKey, durations := range jobs {
		if len(durations) == 0 {
			logE.WithField("job_key", jobKey).Warn("the job is ignored because the job wasn't executed")
			continue
		}
		maxDuration := slices.Max(durations)
		m[jobKey] = int(math.Ceil(maxDuration.Minutes())) + 10 //nolint:mnd
	}

	return m, nil
}
