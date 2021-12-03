package scheduler

import (
	"github.com/jasonlvhit/gocron"
)

// Run new cron job
func Run(run func(), timing string) {
	cronJob := gocron.NewScheduler()
	cronJob.Every(1).Day().At(timing).Do(run)
	<-cronJob.Start()
}
