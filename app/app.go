package app

import (
	"face_management/database"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	Router               *gin.Engine
	KAFKA_CONSUMER_GROUP string

	SQS_QUEUE_TO_EVENT_AGGREGATOR_EVENT string
)

func init() {

	gin.SetMode(gin.ReleaseMode)

	Router = gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isValidCardEnrollmentMode", IsValidCardEnrollmentMode)
		v.RegisterValidation("isDeviceSerialNumber", isDeviceSerialNumber)
		v.RegisterValidation("IsValidCardCommandMsgType", IsValidCardCommandMsgType)
		v.RegisterValidation("IsValidSupportCommandMsgType", IsValidSupportCommandMsgType)

	}

	KAFKA_CONSUMER_GROUP = os.Getenv("KAFKA_CONSUMER_GROUP")

	SQS_QUEUE_TO_EVENT_AGGREGATOR_EVENT = os.Getenv("SQS_QUEUE_TO_EVENT_AGGREGATOR_EVENT")

}

func StartApp() {

	defer func() {
		database.CloseDatabasePool()
	}()

	SetupHealthRoute()
	SetupRoutesMiddleware()
	SetupInternalRoutes()
	SetupIntegratorRoutes()
	SetupSupportRoutes()

	if err := database.InitializeDatabasePool(); err != nil {
		panic(err)
	}

	go func() {

		if err := Router.Run(":" + os.Getenv("SERVER_PORT")); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

}
