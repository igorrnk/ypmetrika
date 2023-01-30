package storage

import (
	"context"
	"encoding/json"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"os"
	"sync"
	"time"
)

type FileStorage struct {
	MemoryStorage
	filename      string
	storeInterval time.Duration
	syncSaveData  bool
	isSaveData    bool
	mutexFile     sync.Mutex
	context       context.Context
}

func NewFileStorage(ctx context.Context, config *configs.ServerConfig) *FileStorage {

	storage := &FileStorage{
		MemoryStorage: MemoryStorage{
			metrics:  make(map[string]*models.Metric, 0),
			mutexMem: sync.RWMutex{},
		},
		context:  ctx,
		filename: config.StoreFileName,
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

func (storage *FileStorage) Write(metric *models.Metric) error {
	err := storage.MemoryStorage.Write(metric)
	if err != nil {
		return err
	}
	if storage.syncSaveData {
		storage.saveData()
	}
	return nil
}

func (storage *FileStorage) tickSave() {
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

func (storage *FileStorage) saveData() {
	storage.mutexFile.Lock()
	defer storage.mutexFile.Unlock()
	file, err := os.Create(storage.filename)
	if err != nil {
		log.Printf("FileStorage.saveData: create error: %v\n", err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("FileStorage.saveData: close error: %v\n", err)
		}
	}()
	storage.MemoryStorage.mutexMem.RLock()
	defer storage.MemoryStorage.mutexMem.RUnlock()
	err = json.NewEncoder(file).Encode(storage.metrics)
	if err != nil {
		log.Printf("FileStorage.saveData: encode error: %v\n", err)
	}
}

func (storage *FileStorage) restoreData() {
	storage.mutexFile.Lock()
	defer storage.mutexFile.Unlock()
	file, err := os.Open(storage.filename)
	if err != nil {
		log.Printf("FileStorage.restoreFromFile: open error: %v\n", err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("FileStorage.restoreFromFile: close error: %v\n", err)
		}
	}()
	decoder := json.NewDecoder(file)
	storage.MemoryStorage.mutexMem.Lock()
	defer storage.MemoryStorage.mutexMem.Unlock()
	err = decoder.Decode(&storage.metrics)
	if err != nil {
		log.Printf("FileStorage.restoreFromFile: decode error: %v\n", err)
	}
}
