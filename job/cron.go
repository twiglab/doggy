package job

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Job interface {
	Run()
}

type JobDefinition interface {
	gocron.JobDefinition
}

type CronWarp struct {
	scheduler gocron.Scheduler
}

func NewCron() (*CronWarp, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &CronWarp{scheduler: s}, nil
}

func (c *CronWarp) AddCronJob(crontab string, job Job) error {
	return c.AddJob(gocron.CronJob(crontab, false), job)
}

func (c *CronWarp) AddDurationJob(d time.Duration, job Job) error {
	return c.AddJob(gocron.DurationJob(d), job)
}

func (c *CronWarp) AddJob(spec JobDefinition, job Job) error {
	return c.AddFunc(spec, JobFunc(job))
}

func (c *CronWarp) AddFunc(spec JobDefinition, cmd func()) error {
	_, err := c.scheduler.NewJob(spec, gocron.NewTask(cmd))
	return err
}

func (c *CronWarp) AddDurationFunc(d time.Duration, cmd func()) error {
	return c.AddFunc(gocron.DurationJob(d), cmd)
}

func (c *CronWarp) AddCronFunc(crontab string, cmd func()) error {
	return c.AddFunc(gocron.CronJob(crontab, false), cmd)
}

func (c *CronWarp) Start() {
	c.scheduler.Start()
}

func (c *CronWarp) Stop() error {
	return c.scheduler.Shutdown()
}

func JobFunc(job Job) func() {
	return func() {
		job.Run()
	}
}
