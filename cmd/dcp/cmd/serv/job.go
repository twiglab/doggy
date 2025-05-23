package serv

import (
	"context"
	"log"

	"github.com/twiglab/doggy/job"
	"github.com/twiglab/doggy/pf"
)

func buildAllJob(ctx context.Context, conf AppConf) *job.CronWarp {
	c, err := job.NewCron()
	if err != nil {
		log.Fatal(err)
	}

	keeplive := buildKeepliveJob(ctx, conf)
	c.AddCronJob(conf.JobConf.Keeplive.Crontab, keeplive)
	return c
}

func buildKeepliveJob(ctx context.Context, conf AppConf) job.Job {
	resolver := ctx.Value(keyCmdb).(pf.DeviceResolver)
	loader := ctx.Value(key_eh).(pf.DeviceLoader)

	return &job.KeepLiveJob{
		DeviceLoader:   loader,
		DeviceResolver: resolver,

		Addr:        conf.JobConf.Keeplive.Addr,
		Port:        conf.JobConf.Keeplive.Port,
		MetadataURL: conf.JobConf.Keeplive.MetadataURL,
	}
}
