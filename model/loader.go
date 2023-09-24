package model

import (
	"fmt"
	"jackpot-mab/reward-predictor/s3"
	"log"
)

const BucketName = "model"

func LoadModels(s3Reader s3.Reader, currentStore Store) Store {
	// Scan bucket a load all files as onnx models.
	modelsToLoad := s3Reader.List(BucketName)

	// Remove non-existent models
	// iterate current store and check if all models exist in bucket
	// if not, unload model.

	for _, m := range modelsToLoad {
		model, err := s3Reader.Read(BucketName, m)
		if err != nil {
			log.Print(fmt.Sprintf("Error reading model from s3: %s", model.Name))
			log.Print(err.Error())
			continue
		}

		err = currentStore.Add(model)
		if err != nil {
			log.Printf("Error loading model: %s", model.Name)
			log.Print(err.Error())
			return currentStore
		}

	}

	return currentStore
}
