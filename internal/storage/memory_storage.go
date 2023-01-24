package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"os"
	"sync"
	"time"
)

type MemoryStorage struct {
	metrics       map[string]*models.Metric
	filename      string
	storeInterval time.Duration
	syncSaveData  bool
	isSaveData    bool
	context       context.Context
	mutexMem      sync.RWMutex
	mutexFile     sync.Mutex
}

func NewAgentStorage() *MemoryStorage {
	storage := &MemoryStorage{
		metrics:  make(map[string]*models.Metric, 0),
		filename: "",
		mutexMem: sync.RWMutex{},
	}
	return storage
}

func NewServerStorage(ctx context.Context, config configs.ServerConfig) *MemoryStorage {
	storage := &MemoryStorage{
		metrics:  make(map[string]*models.Metric, 0),
		context:  ctx,
		filename: config.StoreFileName,
		mutexMem: sync.RWMutex{},
	}
	if config.RestoreData {
		storage.restoreData()
	}
	storage.storeInterval = config.StoreInterval
	if storage.filename != "" {
		if storage.storeInterval == 0 {
			storage.syncSaveData = true
		} else {
			go storage.tickSave()
		}
	}
	return storage
}

func (storage *MemoryStorage) tickSave() {
	tickerSave := time.NewTicker(storage.storeInterval)
	defer tickerSave.Stop()

	for exit := false; !exit; {
		select {
		case <-tickerSave.C:
		case <-storage.context.Done():
			exit = true
		}
		storage.saveData()
	}
}

func (storage *MemoryStorage) Write(metric models.Metric) error {
	storage.mutexMem.Lock()
	defer storage.mutexMem.Unlock()
	storage.metrics[metric.Name] = &metric
	if storage.syncSaveData {
		storage.saveData()
	}
	return nil
}

func (storage *MemoryStorage) Read(metric models.Metric) (models.Metric, error) {
	storage.mutexMem.RLock()
	defer storage.mutexMem.RUnlock()
	if value, ok := storage.metrics[metric.Name]; ok {
		return *value, nil
	}
	return models.Metric{}, fmt.Errorf("metric hasn't been found")
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

func (storage *MemoryStorage) saveData() {
	storage.mutexFile.Lock()
	defer storage.mutexFile.Unlock()
	file, err := os.Create(storage.filename)
	if err != nil {
		log.Printf("MemoryStorage.saveData: create error: %v\n", err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("MemoryStorage.saveData: close error: %v\n", err)
		}
	}()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(storage.metrics)
	if err != nil {
		log.Printf("MemoryStorage.saveData: encode error: %v\n", err)
	}
}

func (storage *MemoryStorage) restoreData() {
	storage.mutexFile.Lock()
	defer storage.mutexFile.Unlock()
	file, err := os.Open(storage.filename)
	if err != nil {
		log.Printf("MemoryStorage.restoreFromFile: open error: %v\n", err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("MemoryStorage.restoreFromFile: close error: %v\n", err)
		}
	}()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storage.metrics)
	if err != nil {
		log.Printf("MemoryStorage.restoreFromFile: decode error: %v\n", err)
	}
}
