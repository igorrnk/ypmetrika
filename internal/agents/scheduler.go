package agents

import (
	"context"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"time"
)

type Scheduler struct {
	UpdateInterval time.Duration
	ReportInterval time.Duration
	Updater        func()
	Reporter       func()
	StopChan       chan time.Time
}

func NewScheduler(conf *configs.AgentConfig, updater func(), reporter func()) *Scheduler {
	newScheduler := &Scheduler{
		conf.PollInterval,
		conf.ReportInterval,
		updater,
		reporter,
		make(chan time.Time),
	}
	return newScheduler
}

func (scheduler Scheduler) Tick(ctx context.Context) {
	tickerPoll := time.NewTicker(scheduler.UpdateInterval)
	tickerReport := time.NewTicker(scheduler.ReportInterval)

OuterLoop:
	for {
		select {
		case <-tickerPoll.C:
			go scheduler.Updater()
		case <-tickerReport.C:
			go scheduler.Reporter()
		case <-ctx.Done():
			break OuterLoop
		}
	}
	tickerPoll.Stop()
	tickerReport.Stop()
}
