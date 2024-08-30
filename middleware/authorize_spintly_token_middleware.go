package middleware

import (
	"encoding/base64"
	"errors"
	"face_management/logger"
	"face_management/misc"
	"face_management/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	SPINTLY_PKEY string
	publicKey    []byte
)

func init() {

	var err error

	SPINTLY_PKEY = os.Getenv("SPINTLY_PKEY")

	publicKey, err = base64.StdEncoding.DecodeString(SPINTLY_PKEY)
	if err != nil {
		logger.Log.Error(err.Error())
		misc.ProcessError(1, err.Error(), "SPINTLY_PKEY base64 decoding failed!")
		panic(err)
	}

}

func AuthorizeSpintlyToken() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			apiError := utils.RenderApiError(http.StatusUnauthorized, 1001, "Authorization token missing in header!")
			c.Abort()
			c.JSON(apiError.StatusCode, apiError.ApplicationError)
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, errors.New("Invalid token. Code 001!")
			}

			key, _ := jwt.ParseRSAPublicKeyFromPEM(publicKey)

			return key, nil

		})

		if err != nil {
			apiError := utils.RenderApiError(http.StatusUnauthorized, 1002, err.Error())
			c.Abort()
			c.JSON(apiError.StatusCode, apiError.ApplicationError)
			return
		}

		claims, claimsIsOk := token.Claims.(jwt.MapClaims)

		if !claimsIsOk {
			apiError := utils.RenderApiError(http.StatusUnauthorized, 1003, "Error: Not able to verfiy token claim")
			c.Abort()
			c.JSON(apiError.StatusCode, apiError.ApplicationError)
			return
		}

		c.Set("caller", 2)
		c.Set("sub", claims["sub"])

	}

}
