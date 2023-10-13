package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ort "github.com/yalue/onnxruntime_go"
	"jackpot-mab/reward-predictor/controller"
	"jackpot-mab/reward-predictor/docs"
	"jackpot-mab/reward-predictor/model"
	"jackpot-mab/reward-predictor/s3"
	"log"
	"net/http"
	"os"
	"strconv"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "jackpot-mab:reward-predictor")
}

func main() {

	ort.SetSharedLibraryPath("third_party/libonnxruntime.1.15.1.dylib")
	err := ort.InitializeEnvironment()
	defer ort.DestroyEnvironment()

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	key := os.Getenv("AWS_KEY")
	secret := os.Getenv("AWS_SECRET")
	endpoint := os.Getenv("ENDPOINT")
	s3ForcePathStyle, err := strconv.ParseBool(os.Getenv("S3_FORCE_PATH_STYLE"))

	if err != nil {
		log.Fatalf("Error loading s3 force path style prop: %v", err)
	}

	s3Reader, err := s3.MakeReader(&s3.AwsConfig{
		Region:           "us-east-1",
		AwsKey:           key,
		AwsSecret:        secret,
		Endpoint:         endpoint,
		S3ForcePathStyle: s3ForcePathStyle,
	})

	if err != nil {
		panic(fmt.Sprintf("Error initializing s3 reader: %s", err.Error()))
	}

	modelStore := model.InitMemoryStore()

	// Cron loader is a process that runs in background that updates the
	// model every n seconds. The Store operations have to be thread safe.
	modelCronLoader := model.InitCronLoader(s3Reader, modelStore)
	go modelCronLoader.Start(120)

	if err != nil {
		panic(fmt.Sprintf("Error initializing model store: %s", err.Error()))
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router := gin.Default()

	predictorController := controller.RewardPredictorController{
		ModelStore: modelStore,
	}

	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/prediction")
		{
			eg.POST("/", predictorController.PredictExperimentRewards)
		}
	}

	router.GET("/", healthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("0.0.0.0:8092")
}
