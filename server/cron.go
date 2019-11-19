package server

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// CronJob 定时任务
func CronJob() {
	log.Info("Start CronJob")

	cron := cron.New()

	//
	cron.AddFunc("0 0 22 * * ?", func() {
		// SendArenaRewardMail(pool.Get(), dbConnect, mqChannel)
	})

	cron.Start()
}
