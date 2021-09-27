package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/lock"
)

var (
	jobs []Job
	c    *cron.Cron

	EmptyOptions []Option

	Recover             = cron.Recover(newLog())
	SkipIfStillRunning  = cron.SkipIfStillRunning(newLog())
	DelayIfStillRunning = cron.DelayIfStillRunning(newLog())
)

type JobWrapper = cron.JobWrapper

type Job interface {
	Name() string
	Spec() string
	Options() []Option
	Run()
}

func Init() error {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	c = cron.New(cron.WithSeconds(), cron.WithLocation(shanghai))

	for _, job := range jobs {
		opt := &Options{}
		for _, o := range job.Options() {
			o(opt)
		}

		var cronJob cron.Job = job

		if opt.SingleNode {
			cronJob = cron.FuncJob(func() {
				l := lock.NewLock(fmt.Sprintf("cron_single_node_%s_%s", config.Service.SrvName, job.Name()))
				if !l.Lock() {
					return
				}
				job.Run()
				l.UnLock()
			})
		}

		if opt.Wrapper != nil {
			cronJob = cron.NewChain(*(opt.Wrapper)).Then(cronJob)
		}

		_, err = c.AddJob(job.Spec(), cronJob)
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
