package agents

import (
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"log"
	"net/http"
	"time"
)

type Config struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	addressServer  string
	nameCSVFile    string
}

var DefaultConfig Config = Config{
	pollInterval:   2 * time.Second,
	reportInterval: 10 * time.Second,
	addressServer:  "127.0.0.1:8080",
	nameCSVFile:    "./configs/metrics.csv",
}

type Agent struct {
	Config  Config
	Metrics *metrics.Metrics
}

func NewAgent() Agent {
	newAgent := Agent{
		Config:  DefaultConfig,
		Metrics: &metrics.Metrics{},
	}
	log.Println("New agent has been made.")
	return newAgent

}

func (agent *Agent) FillMetrics() error {
	err := agent.Metrics.FillFromCSV(agent.Config.nameCSVFile)
	if err != nil {
		log.Printf("FillMetrics(): %v", err)
		return err
	}
	return nil
}

func (agent *Agent) Report() {
	for _, metric := range agent.Metrics.Metrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%s",
			agent.Config.addressServer,
			metric.Type,
			metric.Name,
			metric.Value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Println(err)
			return
		}
		err = resp.Body.Close()
		{
			if err != nil {
				log.Println(err)
			}
		}
		log.Printf("%v Status: %v", url, resp.Status)
	}
	agent.Metrics.UpdateCount = 0
	log.Println("All metrics have been received.")
}

func (agent *Agent) Run() {
	tickerPoll := time.NewTicker(agent.Config.pollInterval)
	tickerReport := time.NewTicker(agent.Config.reportInterval)
	log.Println("Agent is running.")
	for {
		select {
		case <-tickerPoll.C:
			err := agent.Metrics.Update()
			if err != nil {
				log.Println(err)
			}
		case <-tickerReport.C:
			agent.Report()
		}
	}
}
