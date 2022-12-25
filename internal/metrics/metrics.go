package metrics

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type Metric struct {
	Name   string
	Type   string
	Value  string
	Source string
}

type Metrics struct {
	Metrics     []Metric
	UpdateCount int64
}

func (ms *Metrics) Fill() (e error) {
	f, err := os.Open("cmd/agent/metrics.csv")
	defer func(f *os.File) {
		e = f.Close()
	}(f)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	ms.Metrics = make([]Metric, 0, len(records))

	for _, record := range records {
		ms.Metrics = append(ms.Metrics, Metric{
			Name:   record[0],
			Type:   record[1],
			Value:  record[2],
			Source: record[3],
		})
	}
	return nil
}

func (ms *Metrics) Update() error {
	ms.UpdateCount++
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	s := reflect.ValueOf(stats)

	for i, metric := range ms.Metrics {
		switch metric.Source {
		case "runtime":
			ms.Metrics[i].Value = fmt.Sprint(s.FieldByName(metric.Name))
		case "counter":
			ms.Metrics[i].Value = fmt.Sprint(ms.UpdateCount)
		case "random":
			ms.Metrics[i].Value = fmt.Sprint(fmt.Sprint(rand.Int63()))
		}
	}
	log.Println("Metrics have been updated.")
	return nil
}

func (metric *Metric) URLtoMetric(path string) error {
	path = strings.Replace(path, "/", " ", 10)
	n, err := fmt.Sscanf(path, " update %s %s %s\n", &metric.Type, &metric.Name, &metric.Value)
	if err != nil {
		log.Println("URL path is wrong.")
		return err
	}

	if n != 3 {
		log.Println("URL path is wrong.")
		return fmt.Errorf("URLtoMetric")
	}

	return nil
}
