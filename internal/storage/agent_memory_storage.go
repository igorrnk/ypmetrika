package storage

import (
	"encoding/csv"
	"github.com/igorrnk/ypmetrika/internal/models"
	"os"
	"sync"
)

type AgentMemoryStorage struct {
	Metrics       map[string]*models.AgentMetric
	UpdateCounter int64
	Mutex         sync.RWMutex
}

func NewAgentMemoryStorage() *AgentMemoryStorage {
	return &AgentMemoryStorage{
		Metrics:       nil,
		UpdateCounter: 0,
		Mutex:         sync.RWMutex{},
	}
}

func (storage *AgentMemoryStorage) Update(metric models.AgentMetric) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	storage.Metrics[metric.Name].Value = metric.Value
}

func (storage *AgentMemoryStorage) Fill(filename string) (err error) {
	f, err := os.Open(filename)
	defer func(f *os.File) {
		err = f.Close()
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
	storage.Metrics = make(map[string]*models.AgentMetric)

	for _, record := range records {
		storage.Metrics[record[0]] = &models.AgentMetric{
			Name:   record[0],
			Type:   record[1],
			Value:  record[2],
			Source: record[3],
		}
	}
	return nil
}

func (storage *AgentMemoryStorage) ReadAll() []models.AgentMetric {
	storage.Mutex.RLock()
	defer storage.Mutex.RUnlock()
	metrics := make([]models.AgentMetric, 0)
	for _, metric := range storage.Metrics {
		metrics = append(metrics, *metric)
	}
	return metrics
}
