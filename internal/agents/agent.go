package agents

import (
	"context"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/crypts"
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
	Config        *configs.AgentConfig
	Scheduler     *Scheduler
	Repository    models.Repository
	Client        models.Client
	Crypter       models.Crypter
	UpdateCounter int64
}

func NewAgent(config *configs.AgentConfig) (*Agent, error) {

	newAgent := &Agent{
		Config: config,
	}
	newAgent.Scheduler = NewScheduler(config, newAgent.Update, newAgent.Report)
	newAgent.Repository = storage.NewAgentStorage()
	newAgent.Client = delivery.NewRestyClient(config)
	if config.Key != "" {
		newAgent.Crypter = crypts.NewCrypterSHA256(config.Key)
	}

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
				metric.Gauge = float64(field.Uint())
			case reflect.Float64:
				metric.Gauge = field.Float()
			}
		case models.CounterSource:
			metric.Counter = agent.UpdateCounter
		case models.RandomSource:
			metric.Gauge = rand.Float64()
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
		agent.Crypter.AddHash(&metric)
		agent.Client.PostJSON(&metric)
	}
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
