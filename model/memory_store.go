package model

import (
	"errors"
	"io"
	"jackpot-mab/reward-predictor/exchange"
	"sync"
)

type Store interface {
	Get(modelId string) (Instance, error)
	Changed(modelId string, checksum string) bool
	add(modelFile *exchange.File) error
	retainOnly(modelIds []string)
}

type MemoryStore struct {
	mu                 sync.RWMutex
	loadedModelsByName map[string]Instance
}

func InitMemoryStore() Store {
	modelsMap := make(map[string]Instance)
	return &MemoryStore{loadedModelsByName: modelsMap}

}

func (ms *MemoryStore) Get(modelId string) (Instance, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if m, ok := ms.loadedModelsByName[modelId]; ok {
		return m, nil
	}
	return nil, errors.New("model not found")
}

// Changed whether the model inside the Store has the same
// the checksum passed by parameter. Returns true if the model
// does not exist in the store
func (ms *MemoryStore) Changed(modelId string, checksum string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if modelInstanceLoaded, ok := ms.loadedModelsByName[modelId]; ok {
		return modelInstanceLoaded.Checksum() != checksum
	}

	return true
}

func (ms *MemoryStore) retainOnly(modelIds []string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	filtered := make(map[string]Instance)
	for _, m := range modelIds {
		if modelInstanceLoaded, ok := ms.loadedModelsByName[m]; ok {
			filtered[m] = modelInstanceLoaded
		}
	}
	ms.loadedModelsByName = filtered
}

func (ms *MemoryStore) add(modelFile *exchange.File) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	modelBytes, err := io.ReadAll(modelFile.Body)
	if err != nil {
		return err
	}

	model, err := Load(modelBytes, modelFile.Name, *modelFile.Checksum)

	if err != nil {
		return err
	}

	ms.loadedModelsByName[modelFile.Name] = model

	return nil
}
