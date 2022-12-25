package main

import (
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	pollInterval   = 1 * time.Second
	reportInterval = 3 * time.Second
	addressServer  = "127.0.0.1:8080"
)

func Report(ms *metrics.Metrics) {
	log.Println(len(ms.Metrics))
	for _, metric := range ms.Metrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%s",
			addressServer,
			metric.Type,
			metric.Name,
			metric.Value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.Status)
	}
	ms.UpdateCount = 0
	log.Println("All metrics have been receive.")
}

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)
	log.Println("Agent is running.")

	ms := new(metrics.Metrics)
	err := ms.Fill()
	if err != nil {
		log.Fatal(err)
	}

	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	for {
		select {
		case <-tickerPoll.C:
			err := ms.Update()
			if err != nil {
				log.Println(err)
			}
		case <-tickerReport.C:
			Report(ms)
		}
	}
}
