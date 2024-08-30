package utils

import (
	"face_management/models"
	"fmt"
)

const (
	ServiceErrorBaseCode = 0x61000000
)

func RenderApiError(statusCode int, errorCode int, err string) *models.ApiError {

	apiError := models.ApiError{
		StatusCode: statusCode,
		ApplicationError: models.ApplicationError{
			Type: "error",
			Message: models.ApplicationErrorMessage{
				ErrorCode:    errorCode,
				ErrorMessage: fmt.Sprintf("%s", err),
			},
		},
	}

	return &apiError
}

func RenderAppError(errorCode int, err string) *models.ApplicationError {

	appError := models.ApplicationError{
		Type: "error",
		Message: models.ApplicationErrorMessage{
			ErrorCode:    errorCode,
			ErrorMessage: fmt.Sprintf("%s", err),
		},
	}

	return &appError
}

func RenderApiErrorFromAppError(statusCode int, appError *models.ApplicationError) *models.ApiError {

	apiError := models.ApiError{
		StatusCode:       statusCode,
		ApplicationError: *appError,
	}

	return &apiError
}
