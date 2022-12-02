package hook

import "github.com/jasonlvhit/gocron"

type ScheduleItem struct {
	Period *gocron.Job
	Job    func()
}

type Schedule interface {
	Schedulers() []ScheduleItem
}
