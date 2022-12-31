package agents

import (
	"github.com/igorrnk/ypmetrika/configs"
	"time"
)

type Scheduler struct {
	UpdateInterval time.Duration
	ReportInterval time.Duration
	Updater        func()
	Reporter       func()
	StopChan       chan time.Time
}

func NewScheduler(conf configs.AgentConfig, updater func(), reporter func()) *Scheduler {
	newScheduler := &Scheduler{
		conf.PollInterval,
		conf.ReportInterval,
		updater,
		reporter,
		make(chan time.Time),
	}
	return newScheduler
}

func (scheduler Scheduler) Tick() {
	tickerPoll := time.NewTicker(scheduler.UpdateInterval)
	tickerReport := time.NewTicker(scheduler.ReportInterval)

OuterLoop:
	for {
		select {
		case <-tickerPoll.C:
			scheduler.Updater()
		case <-tickerReport.C:
			scheduler.Reporter()
		case <-scheduler.StopChan:
			break OuterLoop

		}
	}
	tickerPoll.Stop()
	tickerReport.Stop()
}

func (scheduler Scheduler) Stop() {
	scheduler.StopChan <- time.Now()
}
