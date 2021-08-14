package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

var (
	jobs []Job
	c    *cron.Cron

	Recover             = cron.Recover(newLog())
	SkipIfStillRunning  = cron.SkipIfStillRunning(newLog())
	DelayIfStillRunning = cron.DelayIfStillRunning(newLog())
)

type JobWrapper = cron.JobWrapper

type Job interface {
	Spec() string
	Wrapper() *JobWrapper
	Run()
}

func Init() error {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	c = cron.New(cron.WithSeconds(), cron.WithLocation(shanghai))

	for _, job := range jobs {
		if job.Wrapper() != nil {
			_, err = c.AddJob(job.Spec(), cron.NewChain(*(job.Wrapper())).Then(job))
		} else {
			_, err = c.AddJob(job.Spec(), job)
		}
		if err != nil {
			return err
		}
	}

	c.Start()

	return nil
}

func AddJobs(jos ...Job) {
	jobs = append(jobs, jos...)
}

func Stop() {
	c.Stop()
}
