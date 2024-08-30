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
	PERMISSION_MANAGEMENT_BASE_URL, PERMISSION_MANAGEMENT_X_API_KEY string
)

func init() {

	PERMISSION_MANAGEMENT_BASE_URL = os.Getenv("PERMISSION_MANAGEMENT_BASE_URL")
	PERMISSION_MANAGEMENT_X_API_KEY = os.Getenv("PERMISSION_MANAGEMENT_X_API_KEY")

}

func GetRequestIdForTransaction(path models.PendingPermissionFaceStatusPathParams, requestBody models.PendingPermissionFaceStatusRequestBody) (*models.PendingPermissionFaceStatusResponse, *models.ApplicationError) {

	headers := map[string]string{"Content-Type": "application/json", "x-api-key": PERMISSION_MANAGEMENT_X_API_KEY}

	url := fmt.Sprintf("%s/permissionManagement/v1/organisations/%d/accessors/%d/faceStatus", PERMISSION_MANAGEMENT_BASE_URL, path.OrganisationId, path.AccessorId)

	statusCode, responseBody, err := clients.RestClient.Post(url, headers, requestBody, time.Duration(time.Duration.Seconds(10)))

	logger.Log.Info("GetRequestIdForTransaction", zap.String("url", url), zap.Any("responseBody", string(responseBody)), zap.Int("statusCode", statusCode))

	if err != nil {
		formattedMessage := fmt.Sprintf("Could not get the request id for this transaction.Error:%s", err.Error())
		logger.Log.Error(formattedMessage)
		appError := utils.RenderAppError(1018, formattedMessage)
		return nil, appError
	}

	if statusCode != 200 {
		appError := utils.RenderAppError(1018, fmt.Sprintf(" Error. Received status code:%d!", statusCode))
		return nil, appError
	}
	var response models.PendingPermissionFaceStatusResponse

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		//TODO add error code
		return nil, utils.RenderAppError(123, err.Error())
	}

	return &response, nil

}
