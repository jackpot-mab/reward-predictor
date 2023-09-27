package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jackpot-mab/reward-predictor/model"
	"log"
	"math/rand"
	"net/http"
)

type RewardPredictorController struct {
	ModelStore model.Store
}

type PredictionRequest struct {
	Model   string    `json:"model"`
	Sample  bool      `json:"sample"`
	Context []float32 `json:"context"`
}

type PredictionResponse struct {
	Prediction float32                `json:"prediction"`
	ExtraData  map[string]interface{} `json:"extra_data"`
}

// PredictExperimentRewards godoc
// @Summary Predict experiment rewards using ONNX models loaded in S3.
// @Description Predicts rewards for a given experiment.
// @ID predict-experiment-rewards
// @Accept json
// @Produce json
// @Param body body PredictionRequest true "Prediction Request"
// @Success 200 {object} PredictionResponse
// @Router /prediction [post]
func (r *RewardPredictorController) PredictExperimentRewards(g *gin.Context) {
	var predictionRequest PredictionRequest
	if err := g.BindJSON(&predictionRequest); err != nil {
		log.Print("error occurred: ", err)
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	modelInstance, err := r.ModelStore.Get(fmt.Sprintf("%s.onnx", predictionRequest.Model))
	if err != nil {
		log.Print("error occurred: ", err)
		g.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prediction, err := modelInstance.Predict(predictionRequest.Context)

	if err != nil {
		log.Print("error occurred", err)
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	if predictionRequest.Sample {
		// sample prediction TODO
		// to test a thompson sampling-like model.
	}

	g.JSON(http.StatusOK, PredictionResponse{
		Prediction: prediction.Label,
		ExtraData: map[string]interface{}{
			// TODO randint for now. but we have to define a way to implement this
			"pulls": rand.Intn(99) + 1,
		},
	})

}
