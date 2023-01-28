package storage

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"sync"
)

type MemoryStorage struct {
	metrics  map[string]*models.Metric
	mutexMem sync.RWMutex
}

func NewAgentStorage() *MemoryStorage {
	storage := &MemoryStorage{
		metrics:  make(map[string]*models.Metric, 0),
		mutexMem: sync.RWMutex{},
	}
	return storage
}

func (storage *MemoryStorage) Write(metric models.Metric) error {
	storage.mutexMem.Lock()
	defer storage.mutexMem.Unlock()
	storage.metrics[metric.Name] = &metric
	return nil
}

func (storage *MemoryStorage) Read(metric models.Metric) (models.Metric, bool) {
	storage.mutexMem.RLock()
	defer storage.mutexMem.RUnlock()
	if value, ok := storage.metrics[metric.Name]; ok {
		return *value, true
	}
	return models.Metric{}, false
}

func (storage *MemoryStorage) ReadAll() ([]models.Metric, error) {
	storage.mutexMem.RLock()
	defer storage.mutexMem.RUnlock()
	metrics := make([]models.Metric, 0)
	for _, metric := range storage.metrics {
		metrics = append(metrics, *metric)
	}
	return metrics, nil
}
