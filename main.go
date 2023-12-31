package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ort "github.com/yalue/onnxruntime_go"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"jackpot-mab/reward-predictor/controller"
	"jackpot-mab/reward-predictor/docs"
	"jackpot-mab/reward-predictor/model"
	"jackpot-mab/reward-predictor/s3"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "jackpot-mab:reward-predictor")
}

func main() {

	shLib := getDefaultSharedLibPath()
	log.Printf("Shared lib onnxruntime: %v", shLib)

	ort.SetSharedLibraryPath(fmt.Sprintf("third_party/%s", shLib))
	err := ort.InitializeEnvironment()
	defer ort.DestroyEnvironment()

	if err != nil {
		log.Printf("Error initializing onnxruntime: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	key := os.Getenv("AWS_KEY")
	secret := os.Getenv("AWS_SECRET")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("ENDPOINT")
	s3ForcePathStyle, err := strconv.ParseBool(os.Getenv("S3_FORCE_PATH_STYLE"))
	s3Bucket := os.Getenv("S3_BUCKET")

	if err != nil {
		log.Fatalf("Error loading s3 force path style prop: %v", err)
	}

	s3Reader, err := s3.MakeReader(&s3.AwsConfig{
		Region:           region,
		AwsKey:           key,
		AwsSecret:        secret,
		Endpoint:         endpoint,
		S3ForcePathStyle: s3ForcePathStyle,
		S3Bucket:         s3Bucket,
	})

	if err != nil {
		panic(fmt.Sprintf("Error initializing s3 reader: %s", err.Error()))
	}

	modelStore := model.InitMemoryStore()

	// Cron loader is a process that runs in background that updates the
	// model every n seconds. The Store operations have to be thread safe.
	modelCronLoader := model.InitCronLoader(s3Reader, modelStore)
	go modelCronLoader.Start(60)

	if err != nil {
		panic(fmt.Sprintf("Error initializing model store: %s", err.Error()))
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

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

func getDefaultSharedLibPath() string {
	// For now, we only include libraries for x86_64 windows, ARM64 darwin, and
	// x86_64 or ARM64 Linux. In the future, libraries may be added or removed.
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "libonnxruntime.1.15.1.dylib"
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return "libonnxruntime.arm.so.1.15.1"
		}
		return "libonnxruntime.intel.so.1.15.1"
	}
	fmt.Printf("Unable to determine a path to the onnxruntime shared library"+
		" for OS \"%s\" and architecture \"%s\".\n", runtime.GOOS,
		runtime.GOARCH)
	return ""
}
