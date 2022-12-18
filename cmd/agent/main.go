package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	addressServer  = "127.0.0.1:8080"
)

type gauge float64
type counter int64

type Metric struct {
	T     string
	Name  string
	Value string
}

var metrics = [len(metricNamesRuntime) + 2]Metric{}
var countUpdate int64 = 0

var metricNamesRuntime = [...]string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"Lookups",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

var metricNamesOther = [...]string{
	"PollCount",
	"RandomValue",
}

func Update() {
	fmt.Println("Update")
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	s := reflect.ValueOf(stats)
	for i, metricName := range metricNamesRuntime {
		v := fmt.Sprint(s.FieldByName(metricName))
		metric := Metric{
			Name:  metricName,
			T:     "gauge",
			Value: v,
		}
		metrics[i] = metric
	}
	i := len(metricNamesRuntime)
	countUpdate++
	metrics[i] = Metric{Name: "PollCount", T: "counter", Value: fmt.Sprint(countUpdate)}
	metrics[i+1] = Metric{Name: "RandomValue", T: "gauge", Value: fmt.Sprint(rand.Int63())}
}

func Report() {
	fmt.Println("Report")
	for _, metric := range metrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%s",
			addressServer,
			metric.T,
			metric.Name,
			metric.Value)
		http.Post(url, "text/plain", nil)
	}
}

func main() {
	//Update()
	//PrintMetrics()
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	for {
		select {
		case <-tickerPoll.C:
			Update()
		case <-tickerReport.C:
			Report()
		}
	}
}
