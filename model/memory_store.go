package model

import (
	"io"
	"jackpot-mab/reward-predictor/exchange"
)

type Store interface {
	Get(modelId string) Instance
	Add(modelFile *exchange.File) error
}

type MemoryStore struct {
	loadedModelsByName map[string]Instance
}

func InitMemoryStore() Store {
	modelsMap := make(map[string]Instance)
	return &MemoryStore{loadedModelsByName: modelsMap}

}

func (ms *MemoryStore) Get(modelId string) Instance {
	return ms.loadedModelsByName[modelId]
}

func (ms *MemoryStore) Add(modelFile *exchange.File) error {

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
