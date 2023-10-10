package model

import (
	ort "github.com/yalue/onnxruntime_go"
	"jackpot-mab/reward-predictor/exchange"
	"log"
)

type Instance interface {
	Predict(features []float32, classes int) (exchange.Prediction, error)
	Checksum() string
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

func (l *LoadedModel) Checksum() string {
	return l.checksum
}

// TODO pensado para un clasificador: porque retorna probabilidades por clase,
// si el modelo es un regresor habría que implementar algún condicional aquí.
func (l *LoadedModel) Predict(features []float32, classes int) (exchange.Prediction, error) {

	inputTensor, e := ort.NewTensor(ort.NewShape(1, int64(len(features))), features)
	defer inputTensor.Destroy()

	outputTensorProba, e := ort.NewEmptyTensor[float32](ort.NewShape(1, int64(classes)))
	outputTensorLabel, e := ort.NewEmptyTensor[int8](ort.NewShape(1))
	defer outputTensorProba.Destroy()
	defer outputTensorLabel.Destroy()

	if e != nil {
		log.Printf("There was an error creating the tensors %v", e)
		return exchange.Prediction{}, e
	}

	e = l.modelSession.Run(
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensorLabel, outputTensorProba})

	if e != nil {
		log.Printf("Error creating the session: %v", e)
		return exchange.Prediction{}, e
	}

	outputProbabilities := outputTensorProba.GetData()
	outputLabel := outputTensorLabel.GetData()

	return exchange.Prediction{
		Label:         float32(outputLabel[0]),
		Probabilities: outputProbabilities,
	}, nil

}
