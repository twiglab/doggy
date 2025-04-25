package job

import "github.com/twiglab/doggy/internal/cron"

type Job interface {
	cron.Job
}

type CronWarp struct {
	cron *cron.Cron
}

func NewCron() *CronWarp {
	return &CronWarp{
		cron: cron.New(),
	}
}

func (c *CronWarp) AddJob(spec string, job Job) error {
	_, err := c.cron.AddJob(spec, job)
	return err
}

func (c *CronWarp) Start() {
	c.cron.Start()
}

func (c *CronWarp) Stop() {
	c.cron.Stop()
}
