package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jackpot-mab/reward-predictor/controller"
	"jackpot-mab/reward-predictor/docs"
	"log"
	"net/http"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "jackpot-mab:reward-predictor")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// key := os.Getenv("AWS_KEY")
	// secret := os.Getenv("AWS_SECRET")
	// aws_secret := os.Getenv("AWS_SECRET")

	docs.SwaggerInfo.BasePath = "/api/v1"
	router := gin.Default()

	predictorController := controller.RewardPredictorController{}

	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/prediction")
		{
			eg.GET("/", predictorController.PredictExperimentRewards)
		}
	}

	router.GET("/", healthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("localhost:8092")
}
