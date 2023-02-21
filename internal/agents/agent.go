package agents

import (
	"context"
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/crypts"
	"github.com/igorrnk/ypmetrika/internal/delivery"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"time"
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
	newAgent.Scheduler = NewScheduler(config, newAgent.Update, newAgent.ReportBatch)
	newAgent.Repository = storage.NewAgentStorage()
	newAgent.Client = delivery.NewRestyClient(config)
	newAgent.Crypter = crypts.NewCrypterSHA256(config.Key)

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
	go agent.UpdateMain()
	go agent.UpdateAdd()

}

func (agent *Agent) UpdateMain() {
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
		default:
			continue
		}
		err := agent.Repository.Write(metric)
		if err != nil {
			log.Println(err)
		}
		//log.Printf("Metric %v (%v) = %v has been updated.", metric.Name, metric.Type, metric.Value)
	}
	log.Println("Metrics have been updated.")
}

func (agent *Agent) UpdateAdd() {
	v, _ := mem.VirtualMemory()
	n, _ := cpu.Counts(true)
	p, _ := cpu.Percent(time.Second, true)
	for _, metric := range models.AllMetrics {
		switch metric.Source {
		case models.GopsutilSource:
			if metric.Name == "TotalMemory" {
				metric.Gauge = float64(v.Total)
			}
			if metric.Name == "FreeMemory" {
				metric.Gauge = float64(v.Free)
			}
			if metric.Name == "CPUutilization" {
				for i := 0; i < n; i++ {
					metric1 := &models.Metric{
						Name:  metric.Name + fmt.Sprint(i+1),
						Type:  models.GaugeType,
						Gauge: p[i],
					}
					err := agent.Repository.Write(metric1)
					if err != nil {
						log.Println(err)
					}
				}
				continue
			}
		default:
			continue
		}
		err := agent.Repository.Write(metric)
		if err != nil {
			log.Println(err)
		}
		//log.Printf("Metric %v (%v) = %v has been updated.", metric.Name, metric.Type, metric.Value)
	}

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

func (agent *Agent) ReportBatch() {
	metrics, err := agent.Repository.ReadAll()
	if err != nil {
		log.Printf("agents.Report: reporting: %v", err)
		return
	}

	for i := range metrics {
		agent.Crypter.AddHash(&metrics[i])
	}
	agent.Client.PostMetrics(metrics)
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
