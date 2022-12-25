package storage

import (
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"log"
	"strconv"
)

type Repositories interface {
	Write(m *metrics.Metric) error
}

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func (memStorage *MemStorage) Write(metric *metrics.Metric) error {
	if memStorage.GaugeMetrics == nil {
		memStorage.GaugeMetrics = make(map[string]float64, 1)
	}
	if memStorage.CounterMetrics == nil {
		memStorage.CounterMetrics = make(map[string]int64, 1)
	}
	var err error
	switch metric.Type {
	case "gauge":
		err = memStorage.AddGaugeMetric(metric.Name, metric.Value)
	case "counter":
		err = memStorage.AddCounterMetric(metric.Name, metric.Value)
	}
	if err != nil {
		return err
	}
	log.Printf("Metric name: %v, type: %v, value: %v is added\n", metric.Name, metric.Type, metric.Value)
	return nil
}

func (memStorage *MemStorage) AddGaugeMetric(name string, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	memStorage.GaugeMetrics[name] = v
	return nil
}

func (memStorage *MemStorage) AddCounterMetric(name string, value string) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	memStorage.CounterMetrics[name] = memStorage.CounterMetrics[name] + v
	return nil
}
