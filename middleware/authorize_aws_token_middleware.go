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
	AWS_KID1, AWS_PKEY1, AWS_KID2, AWS_PKEY2 string
	awsPublicKey1, awsPublicKey2             []byte
)

func init() {

	var err error

	AWS_KID1 = os.Getenv("AWS_KID1")
	AWS_PKEY1 = os.Getenv("AWS_PKEY1")
	AWS_KID2 = os.Getenv("AWS_KID2")
	AWS_PKEY2 = os.Getenv("AWS_PKEY2")

	awsPublicKey1, err = base64.StdEncoding.DecodeString(AWS_PKEY1)
	if err != nil {
		logger.Log.Error(err.Error())
		misc.ProcessError(1, err.Error(), "AWS_PKEY1 base64 decoding failed!")
		panic(err)
	}

	awsPublicKey2, err = base64.StdEncoding.DecodeString(AWS_PKEY2)
	if err != nil {
		logger.Log.Error(err.Error())
		misc.ProcessError(1, err.Error(), "AWS_PKEY2 base64 decoding failed!")
		panic(err)
	}

}

func AuthorizeAwsToken() gin.HandlerFunc {

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

			keyId, keyFound := token.Header["kid"]

			if !keyFound {
				return nil, errors.New("Invalid token. Code 002!")
			}

			keyIdStr, keyIdIsStr := keyId.(string)

			if !keyIdIsStr {
				return nil, errors.New("Invalid token. Code 003!")
			}

			var publicKey []byte

			if keyIdStr == AWS_KID1 {
				publicKey = awsPublicKey1
			} else if keyIdStr == AWS_KID2 {
				publicKey = awsPublicKey2
			} else {
				return nil, errors.New("Invalid token. Code 004!")
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

			return key, err

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

		c.Set("caller", 1)
		c.Set("scopes", claims["custom:userScopes"])
		c.Set("phoneNumber", claims["phone_number"])

	}

}
