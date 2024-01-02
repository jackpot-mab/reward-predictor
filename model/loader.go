package model

import (
	"fmt"
	"jackpot-mab/reward-predictor/metrics"
	"jackpot-mab/reward-predictor/s3"
	"log"
	"time"
)

const BucketName = "jackpot-bucket"

// CronLoader is a process that is intended to run in the
// background that reads the model S3 folder and if it
// detects some model file has changed, it triggers an update
// of the model in the model.Store. Also removes from repo
// models erased from S3.
type CronLoader interface {
	Start(reloadTimeSeconds int)
	Stop()
}

type CronLoaderImpl struct {
	stopChan   chan struct{}
	s3Reader   s3.Reader
	modelStore Store
}

func InitCronLoader(s3Reader s3.Reader, initialStore Store) CronLoader {
	return &CronLoaderImpl{
		s3Reader:   s3Reader,
		modelStore: initialStore,
	}
}

func (c *CronLoaderImpl) loadModels() {

	// Scan bucket and load all files as onnx models.
	modelsToLoad := c.s3Reader.List(BucketName)

	for _, m := range modelsToLoad {
		model, err := c.s3Reader.Read(BucketName, m)
		if err != nil {
			log.Print(fmt.Sprintf("Error reading model from s3: %s", model.Name))
			log.Print(err.Error())
			continue
		}

		// TODO perhaps the checksum (aka ETag) could be download
		// first, before the file, and check the existence matching the Etag with saved checksum.
		if c.modelStore.Changed(model.Name, *model.Checksum) {
			err = c.modelStore.add(model)
			if err != nil {
				log.Printf("Error loading model: %s", model.Name)
				log.Print(err.Error())
				return
			}
			metrics.ModelUpdated.WithLabelValues(model.Name).Inc()
		}

		model.Body.Close()

	}

	c.modelStore.retainOnly(modelsToLoad)

	return
}

func (c *CronLoaderImpl) Start(reloadTimeSeconds int) {
	c.stopChan = make(chan struct{})
	ticker := time.NewTicker(time.Duration(reloadTimeSeconds) * time.Second)
	defer ticker.Stop()

	c.loadModels()
	for {
		select {
		case <-ticker.C:
			c.loadModels()

		case <-c.stopChan:
			fmt.Println("Stopping background task.")
			return
		}
	}

}

func (c *CronLoaderImpl) Stop() {
	<-c.stopChan
}
