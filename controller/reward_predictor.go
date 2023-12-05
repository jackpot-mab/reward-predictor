package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jackpot-mab/reward-predictor/metrics"
	"jackpot-mab/reward-predictor/model"
	"log"
	"net/http"
)

type RewardPredictorController struct {
	ModelStore model.Store
}

type PredictionRequest struct {
	Model   string    `json:"model"`
	Sample  bool      `json:"sample"`
	Context []float32 `json:"context"`
	Classes []string  `json:"classes"`
}

type PredictionResponse struct {
	Prediction float32 `json:"prediction"`
}

// PredictExperimentRewards godoc
// @Summary Predict experiment rewards using ONNX models loaded in S3.
// @Description Predicts rewards for a given experiment.
// @ID predict-experiment-rewards
// @Accept json
// @Produce json
// @Param body PredictionRequest true "Prediction Request"
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

	prediction, err := modelInstance.Predict(predictionRequest.Context, len(predictionRequest.Classes))

	if err != nil {
		log.Print("error occurred", err)
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if predictionRequest.Sample {
		// sample prediction TODO
		// to test a thompson sampling-like model.
	}

	metrics.ModelPredictions.WithLabelValues(predictionRequest.Model).Inc()

	g.JSON(http.StatusOK, PredictionResponse{
		Prediction: prediction.Label,
	})

}
