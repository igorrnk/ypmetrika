package agents

import (
	"fmt"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/delivery"
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
	Repository    storage.AgentRepository
	Client        delivery.Client
	UpdateCounter int64
}

func NewAgent(config configs.AgentConfig) (*Agent, error) {

	newAgent := &Agent{
		Config: config,
	}
	newAgent.Scheduler = NewScheduler(config, newAgent.Update, newAgent.Report)
	newAgent.Repository = storage.NewAgentMemoryStorage()
	newAgent.Repository.Fill(config.NameCSVFile)
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

	for _, metric := range agent.Repository.ReadAll() {
		var value string
		switch metric.Source {
		case "runtime":
			value = fmt.Sprint(s.FieldByName(metric.Name))
		case "counter":
			value = fmt.Sprint(agent.UpdateCounter)
		case "random":
			value = fmt.Sprint(fmt.Sprint(rand.Int63()))
		}
		metric.Value = value
		agent.Repository.Update(metric)
	}
	log.Println("Metrics have been updated.")
	agent.UpdateCounter++
}

func (agent *Agent) Report() {
	for _, metric := range agent.Repository.ReadAll() {
		agent.Client.Post(&metric)
	}
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
