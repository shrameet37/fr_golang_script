package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"face_management/utils"
)

type API_KEY_TYPE int

const (
	API_KEY API_KEY_TYPE = iota
	SUPPORT_API_KEY
)

var (
	apiKeys = map[API_KEY_TYPE]string{
		API_KEY:         os.Getenv("API_KEY"),
		SUPPORT_API_KEY: os.Getenv("SUPPORT_API_KEY"),
	}
)

func AuthorizeApiKey(keyType API_KEY_TYPE) gin.HandlerFunc {
	return func(c *gin.Context) {

		xApiKey := c.GetHeader("x-api-key")

		if xApiKey == "" {
			HandleUnauthorized(c, http.StatusUnauthorized, 1001, "X-API-KEY missing in header!")
			return
		}

		expectedApiKey, ok := apiKeys[keyType]
		if !ok || xApiKey != expectedApiKey {
			HandleUnauthorized(c, http.StatusUnauthorized, 1002, "API key mismatch!")
			return
		}

		c.Set("caller", 3)
		c.Set("sub", "internal")
	}
}

func HandleUnauthorized(c *gin.Context, statusCode int, errorCode int, message string) {

	apiError := utils.RenderApiError(statusCode, errorCode, message)
	c.Abort()
	c.JSON(apiError.StatusCode, apiError.ApplicationError)
}
