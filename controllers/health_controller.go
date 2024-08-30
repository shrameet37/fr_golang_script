package controllers

import (
	"face_management/models"
	"face_management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HEALTH_CONTROLLER_BYTE3_VALUE = 0x00100000 //1 = controller, 0 = health
)

var (
	HEALTH_CONTROLLER_ERROR_BASE_CODE int
)

func init() {
	HEALTH_CONTROLLER_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + HEALTH_CONTROLLER_BYTE3_VALUE
}

func GetHealth(c *gin.Context) {

	c.JSON(http.StatusOK, models.SuccessResponse{Type: "success", Message: "v3: service is healthy!!!"})

}

func NoRoute(c *gin.Context) {

	c.JSON(http.StatusNotFound, models.ApplicationError{Type: "error", Message: models.ApplicationErrorMessage{ErrorCode: HEALTH_CONTROLLER_ERROR_BASE_CODE + 1, ErrorMessage: "Route not found"}})

}
