package misc

import (
	"face_management/database"
	"face_management/logger"
	"face_management/models"
	"fmt"
)

const (
	MISC_ERROR                       = 1
	KAFKA_ERROR_NO_ROLLBACK_REQUIRED = 2
	API_ERROR_NO_ROLLBACK_REQUIRED   = 3
	ERROR_NO_ROLLBACK_REQUIRED       = 4
	ERROR_REQUIRE_ROLLBACK           = 5
	KAFKA_ERROR_REQUIRE_ROLLBACK     = 6
	API_ERROR_REQUIRE_ROLLBACK       = 7
	KAFKA_PRODUCER_ERROR             = 8
	KAFKA_CONSUMER_ERROR             = 9
)

func ProcessError(priority int, errorMessage string, additionalInfo interface{}) {

	switch value := additionalInfo.(type) {

	case []byte:

		additionalInfoConverted, ok := additionalInfo.([]byte)

		if ok {

			appError := database.ErDb.ProcessErrorMessages(priority, errorMessage, string(additionalInfoConverted))

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as byte array could not be converted to string.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case *models.ApplicationError:

		additionalInfoConverted, ok := additionalInfo.(*models.ApplicationError)

		if ok {

			additionalInfoToWriteToDb := fmt.Sprintf("ApplicationErrorMessage Details-ErrorCode:%d,ErrorMessage:%s", additionalInfoConverted.Message.ErrorCode, additionalInfoConverted.Message.ErrorMessage)

			appError := database.ErDb.ProcessErrorMessages(priority, errorMessage, additionalInfoToWriteToDb)

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as *models.ApplicationError message could not be converted.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case *models.ApiError:

		additionalInfoConverted, ok := additionalInfo.(*models.ApiError)

		if ok {

			additionalInfoToWriteToDb := fmt.Sprintf("ApiErrorMessage Details-StatusCode:%d,ErrorCode:%d,ErrorMessage:%s", additionalInfoConverted.StatusCode, additionalInfoConverted.ApplicationError.Message.ErrorCode, additionalInfoConverted.ApplicationError.Message.ErrorMessage)

			appError := database.ErDb.ProcessErrorMessages(priority, errorMessage, additionalInfoToWriteToDb)

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as *models.ApplicationError message could not be converted.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case nil:

		appError := database.ErDb.ProcessErrorMessages(priority, errorMessage, "")

		if appError != nil {

			errorMsg := fmt.Sprintf("Failed to write ProcessError message to database for nil message.Message priority:%d,Message error:%s", priority, errorMessage)
			logger.Log.Error(errorMsg)
		}

	default:
		errorMsg := fmt.Sprintf("Could not process ProcessError message as type of data interface is not supported.Data interface type:%v,Message priority:%d,Message error:%s,Message data:%v", value, priority, errorMessage, additionalInfo)
		logger.Log.Error(errorMsg)

	}

}
