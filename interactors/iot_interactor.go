package interactors

import (
	"face_management/clients"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	IOT_ENGINE_X_API_KEY, IOT_ENGINE_BASE_URL string
)

func init() {

	ENV = os.Getenv("ENV")
	IOT_ENGINE_BASE_URL = os.Getenv("IOT_ENGINE_BASE_URL")
	IOT_ENGINE_X_API_KEY = os.Getenv("IOT_ENGINE_X_API_KEY")

}

func SendMsg(requestBody interface{}) *models.ApplicationError {

	headers := map[string]string{"Content-Type": "application/json", "x-api-key": IOT_ENGINE_X_API_KEY}

	url := fmt.Sprintf("%s/v2/message/receiver", IOT_ENGINE_BASE_URL)

	statusCode, responseBody, err := clients.RestClient.Post(url, headers, requestBody, time.Duration(time.Duration.Seconds(15)))

	logger.Log.Info("GetKeyIdFromDataRepo", zap.String("url", url), zap.Any("responseBody", string(responseBody)), zap.Int("statusCode", statusCode))

	if err != nil {
		formattedMessage := fmt.Sprintf("Could not get the key id for this transaction.Error:%s", err.Error())
		logger.Log.Error(formattedMessage)
		appError := utils.RenderAppError(1018, formattedMessage)
		return appError
	}

	if statusCode != 200 {
		appError := utils.RenderAppError(1018, fmt.Sprintf(" Error. Received status code:%d!", statusCode))
		return appError
	}
	return nil

}
