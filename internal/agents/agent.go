package agents

import (
	"github.com/igorrnk/ypmetrika/configs"
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
	newAgent.Repository = storage.NewMemoryStorage()
	newAgent.Client = delivery.NewRestyClient(config)

	return newAgent, nil
}

func (agent *Agent) Run() error {
	log.Println("Agent is running.")
	go agent.Scheduler.Tick()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	agent.Scheduler.Stop()

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
				metric.Value.Gauge = float64(field.Uint())
			case reflect.Float64:
				metric.Value.Gauge = field.Float()
			}
		case models.CounterSource:
			metric.Value.Counter = agent.UpdateCounter
		case models.RandomSource:
			metric.Value.Gauge = rand.Float64()
		}
		err := agent.Repository.Write(metric)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Metric %v (%v) = %v has been updated.", metric.Name, metric.Type, metric.Value)
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
		agent.Client.Post(&metric)
	}
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
