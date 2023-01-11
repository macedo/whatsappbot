package scheduler

import (
	"time"

	"github.com/procyon-projects/chrono"
)

var scheduler chrono.TaskScheduler

func ScheduleWithTime(task chrono.Task, when time.Time) (chrono.ScheduledTask, error) {
	return scheduler.Schedule(task, chrono.WithTime(when))
}

func Start() {
	scheduler = chrono.NewDefaultTaskScheduler()
}

func Shutdown() {
	scheduler.Shutdown()
}
