package agents

import (
	"context"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/delivery"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
)

type Agent struct {
	Config        configs.AgentConfig
	Scheduler     *Scheduler
	Repository    models.Repository
	Client        models.Client
	UpdateCounter int64
}

func NewAgent(config configs.AgentConfig) (*Agent, error) {

	newAgent := &Agent{
		Config: config,
	}
	newAgent.Scheduler = NewScheduler(config, newAgent.Update, newAgent.Report)
	newAgent.Repository = storage.NewAgentStorage()
	newAgent.Client = delivery.NewRestyClient(config)

	return newAgent, nil
}

func (agent *Agent) Run() error {
	log.Println("Agent is running.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go agent.Scheduler.Tick(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Agent has been stopped.")
	return nil
}

func (agent *Agent) Update() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	s := reflect.ValueOf(stats)
	agent.UpdateCounter++
	for _, metric := range models.AllMetrics {
		switch metric.Source {
		case models.RuntimeSource:
			field := s.FieldByName(metric.Name)
			switch field.Kind() {
			case reflect.Uint64, reflect.Uint32:
				u := float64(field.Uint())
				metric.Value.Gauge = &u
			case reflect.Float64:
				u := field.Float()
				metric.Value.Gauge = &u
			}
		case models.CounterSource:
			c := agent.UpdateCounter
			metric.Value.Counter = &c
		case models.RandomSource:
			f := rand.Float64()
			metric.Value.Gauge = &f
		}
		err := agent.Repository.Write(metric)
		if err != nil {
			log.Println(err)
		}
		//log.Printf("Metric %v (%v) = %v has been updated.", metric.Name, metric.Type, metric.Value)
	}
	log.Println("Metrics have been updated.")
}

func (agent *Agent) Report() {
	metrics, err := agent.Repository.ReadAll()
	if err != nil {
		log.Printf("agents.Report: reporting: %v", err)
		return
	}
	for _, metric := range metrics {
		agent.Client.PostJSON(&metric)
	}
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
