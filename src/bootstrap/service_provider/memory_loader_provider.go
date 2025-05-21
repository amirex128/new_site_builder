package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"git.snappfood.ir/backend/go/packages/sf-memory-loader"
	globalloader "github.com/amirex128/new_site_builder/src/internal/infra/memory_loader"
	"time"
)

func MemoryLoaderProvider(logger sflogger.Logger) {

	go sfmemoryloader.NewScheduler(logger).Duration(
		sfmemoryloader.DurationManager{
			Handler:  globalloader.DealProjectLoader,
			Duration: time.Minute * 5,
		},
	).CronJob(
		sfmemoryloader.CronJobManager{
			Handler: globalloader.DealProjectLoader,
			Crontab: "0 0 * * *", // Run at midnight
		},
	).Daily(
		sfmemoryloader.DailyManager{
			Handler:  globalloader.DealProjectLoader,
			Interval: 1,
			AtTimes: []time.Time{
				time.Date(0, 0, 0, 8, 0, 0, 0, time.Local),
				time.Date(0, 0, 0, 17, 0, 0, 0, time.Local),
			},
		},
	).Start()

}
