package storage

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"sync"
)

type MemoryStorage struct {
	Metrics map[string]*models.Metric
	Mutex   sync.RWMutex
}

func New() *MemoryStorage {
	return &MemoryStorage{
		Metrics: make(map[string]*models.Metric, 0),
		Mutex:   sync.RWMutex{},
	}
}

func (storage *MemoryStorage) Write(metric models.Metric) error {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	storage.Metrics[metric.Name] = &metric
	return nil
}

func (storage *MemoryStorage) Read(metric models.Metric) (models.Metric, bool) {
	storage.Mutex.RLock()
	defer storage.Mutex.RUnlock()
	if value, ok := storage.Metrics[metric.Name]; ok {
		return *value, true
	}
	return models.Metric{}, false
}

func (storage *MemoryStorage) ReadAll() ([]models.Metric, error) {
	storage.Mutex.RLock()
	defer storage.Mutex.RUnlock()
	metrics := make([]models.Metric, 0)
	for _, metric := range storage.Metrics {
		metrics = append(metrics, *metric)
	}
	return metrics, nil
}
