package agents

import (
	"context"
	"errors"
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
	"sync/atomic"
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
	//agent.UpdateCounter++
	atomic.AddInt64(&agent.UpdateCounter, 1)
	for _, metric := range models.AllMetrics {
		m := &models.Metric{
			Name: metric.Name,
			Type: metric.Type,
		}
		switch metric.Source {
		case models.RuntimeSource:
			field := s.FieldByName(metric.Name)
			switch field.Kind() {
			case reflect.Uint64, reflect.Uint32:
				m.Gauge = float64(field.Uint())
			case reflect.Float64:
				m.Gauge = field.Float()
			}
		case models.CounterSource:
			m.Counter = agent.UpdateCounter
		case models.RandomSource:
			m.Gauge = rand.Float64()
		default:
			continue
		}
		err := agent.Repository.Write(m)
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
		m := &models.Metric{
			Name: metric.Name,
			Type: metric.Type,
		}
		if metric.Source == models.GopsutilSource {
			if metric.Name == "TotalMemory" {
				m.Gauge = float64(v.Total)
			}
			if metric.Name == "FreeMemory" {
				m.Gauge = float64(v.Free)
			}
			if metric.Name == "CPUutilization" {
				for i := 0; i < n; i++ {
					m1 := &models.Metric{
						Name:  metric.Name + fmt.Sprint(i+1),
						Type:  models.GaugeType,
						Gauge: p[i],
					}
					err := agent.Repository.Write(m1)
					if err != nil {
						log.Println(err)
					}
				}
				continue
			}
		} else {
			continue
		}
		err := agent.Repository.Write(m)
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
		err = agent.Client.PostJSON(&metric)
		if errors.Is(err, models.ErrNotReport) {
			log.Printf("agents.Report: reporting: %v", err)
		}
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
	//agent.Client.PostMetrics(metrics)
	err = agent.Client.PostMetrics(metrics)
	if errors.Is(err, models.ErrNotReport) {
		log.Printf("agents.Report: reporting: %v", err)
	}
	agent.UpdateCounter = 0
	log.Println("Metrics have been posted.")
}
