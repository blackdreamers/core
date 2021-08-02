package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

var (
	jobs []Job
	c    *cron.Cron
)

type Job interface {
	Spec() string
	Run()
}

func Init() error {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	c = cron.New(cron.WithSeconds(), cron.WithLocation(shanghai))

	for _, job := range jobs {
		_, err = c.AddJob(job.Spec(), job)
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

func DelayIfStillRunJob(job cron.Job) cron.Job {
	return cron.NewChain(cron.DelayIfStillRunning(newLog())).Then(job)
}

func SkipIfStillRunJob(job cron.Job) cron.Job {
	return cron.NewChain(cron.SkipIfStillRunning(newLog())).Then(job)
}
