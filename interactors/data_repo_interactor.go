package interactors

import (
	"encoding/json"
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
	ENV                                                 string
	DATA_REPOSITORY_BASE_URL, DATA_REPOSITORY_X_API_KEY string
)

func init() {

	ENV = os.Getenv("ENV")
	DATA_REPOSITORY_BASE_URL = os.Getenv("DATA_REPOSITORY_URL")
	DATA_REPOSITORY_X_API_KEY = os.Getenv("DATA_REPOSITORY_API_KEY")

}

func AddDataToDataRepo(requestBody *models.AddDataToDataRepoRequest) (*models.AddDataToDataRepoResponse, *models.ApplicationError) {

	headers := map[string]string{"Content-Type": "application/json", "x-api-key": DATA_REPOSITORY_X_API_KEY}

	url := fmt.Sprintf("%s/v1/deviceKeys", DATA_REPOSITORY_BASE_URL)

	statusCode, responseBody, err := clients.RestClient.Post(url, headers, requestBody, time.Duration(time.Duration.Seconds(15)))

	logger.Log.Info("GetKeyIdFromDataRepo", zap.String("url", url), zap.Any("responseBody", string(responseBody)), zap.Int("statusCode", statusCode))

	if err != nil {
		formattedMessage := fmt.Sprintf("Could not get the key id for this transaction.Error:%s", err.Error())
		logger.Log.Error(formattedMessage)
		appError := utils.RenderAppError(1018, formattedMessage)
		return nil, appError
	}

	if statusCode != 200 {
		appError := utils.RenderAppError(1018, fmt.Sprintf(" Error. Received status code:%d!", statusCode))
		return nil, appError
	}
	var response models.AddDataToDataRepoResponse

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		//TODO add error code
		return nil, utils.RenderAppError(123, err.Error())
	}

	return &response, nil

}

func DeleteDataRepoKeyId(keyId int) (bool, *models.ApplicationError) {

	headers := map[string]string{"Content-Type": "application/json", "x-api-key": DATA_REPOSITORY_X_API_KEY}

	url := fmt.Sprintf("%s/v1/deviceKeys/%d", DATA_REPOSITORY_BASE_URL, keyId)

	statusCode, responseBody, err := clients.RestClient.Delete(url, headers, time.Duration(time.Duration.Seconds(10)))

	logger.Log.Info("GetKeyIdFromDataRepo", zap.String("url", url), zap.Any("responseBody", string(responseBody)), zap.Int("statusCode", statusCode))

	if err != nil {
		formattedMessage := fmt.Sprintf("Could not get the key id for this transaction.Error:%s", err.Error())
		logger.Log.Error(formattedMessage)
		appError := utils.RenderAppError(1018, formattedMessage)
		return false, appError
	}

	if statusCode != 200 {
		appError := utils.RenderAppError(1018, fmt.Sprintf(" Error. Received status code:%d!", statusCode))
		return false, appError
	}
	var response map[string]interface{}
	json.Unmarshal(responseBody, &response)

	rtype, ok := response["type"].(string)
	if !ok {
		logger.Log.Error("DeleteDataRepoKeyId", zap.Any("key missing", "type"))
		appError := utils.RenderAppError(1018, "key missing in response")
		return false, appError
	}

	if rtype == "success" {
		return true, nil
	} else {
		logger.Log.Error("DeleteDataRepoKeyId", zap.Any("API failed", string(responseBody)))
		appError := utils.RenderAppError(1018, "API failed")
		return false, appError
	}
}
