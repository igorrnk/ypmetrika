package storage

import (
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"strconv"
	"sync"
)

type ServerMemoryStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
	Mutex          sync.RWMutex
}

func NewServerMemoryStorage() *ServerMemoryStorage {
	return &ServerMemoryStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64)}
}

func (memStorage *ServerMemoryStorage) Write(metric models.ServerMetric) error {
	memStorage.Mutex.Lock()
	defer memStorage.Mutex.Unlock()
	var err error
	switch metric.Type {
	case "gauge":
		err = memStorage.addGaugeMetric(metric.Name, metric.Value)
	case "counter":
		err = memStorage.addCounterMetric(metric.Name, metric.Value)
	}
	if err != nil {
		return err
	}
	log.Printf("AgentMetric name: %v, type: %v, value: %v is added\n", metric.Name, metric.Type, metric.Value)
	return nil
}

func (memStorage *ServerMemoryStorage) Read(metric models.ServerMetric) (models.ServerMetric, bool) {
	memStorage.Mutex.RLock()
	defer memStorage.Mutex.RUnlock()
	switch metric.Type {
	case "gauge":
		if value, ok := memStorage.GaugeMetrics[metric.Name]; ok {
			metric.Value = fmt.Sprint(value)
			return metric, true
		}
	case "counter":
		if value, ok := memStorage.CounterMetrics[metric.Name]; ok {
			metric.Value = fmt.Sprint(value)
			return metric, true
		}
	}
	return models.ServerMetric{}, false
}

func (memStorage *ServerMemoryStorage) ReadAll() []models.ServerMetric {
	memStorage.Mutex.RLock()
	defer memStorage.Mutex.RUnlock()
	metrics := make([]models.ServerMetric, 0)
	for key, value := range memStorage.GaugeMetrics {
		metrics = append(metrics, models.ServerMetric{
			Name:  key,
			Value: fmt.Sprint(value),
		})
	}
	for key, value := range memStorage.CounterMetrics {
		metrics = append(metrics, models.ServerMetric{
			Name:  key,
			Value: fmt.Sprint(value),
		})
	}
	return metrics
}

func (memStorage *ServerMemoryStorage) addGaugeMetric(name string, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	memStorage.GaugeMetrics[name] = v
	return nil
}

func (memStorage *ServerMemoryStorage) addCounterMetric(name string, value string) error {
	//TODO логику перенести в server
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	memStorage.CounterMetrics[name] = memStorage.CounterMetrics[name] + v
	return nil
}
