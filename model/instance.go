package model

import (
	ort "github.com/yalue/onnxruntime_go"
	"log"
)

type Instance interface {
	Predict(features []float32) int
}

type LoadedModel struct {
	name         string
	checksum     string
	modelSession *ort.DynamicAdvancedSession
}

func Load(modelBytes []byte, name string, checksum string) (Instance, error) {

	session, e := ort.NewDynamicAdvancedSessionWithONNXData(modelBytes,
		[]string{"input"},
		[]string{"label", "probabilities"},
		nil)

	if e != nil {
		return nil, e
	}

	return &LoadedModel{
		checksum:     checksum,
		name:         name,
		modelSession: session,
	}, nil
}

func (l *LoadedModel) Predict(features []float32) int {

	inputData := features

	inputTensor, e := ort.NewTensor(ort.NewShape(1, 8), inputData)
	defer inputTensor.Destroy()

	outputTensorProba, e := ort.NewEmptyTensor[float32](ort.NewShape(1, 2))
	outputTensorLabel, e := ort.NewEmptyTensor[int8](ort.NewShape(1))
	defer outputTensorProba.Destroy()
	defer outputTensorLabel.Destroy()

	if e != nil {
		log.Printf("Erro %v", e)
	}

	e = l.modelSession.Run(
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensorLabel, outputTensorProba})

	if e != nil {
		log.Printf("Error creating the session: %v", e)
	}

	outputData := outputTensorProba.GetData()
	print(outputData)

	outpu2 := outputTensorLabel.GetData()
	print(outpu2)

	return int(outpu2[0])

}
