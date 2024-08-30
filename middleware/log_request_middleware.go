package middleware

import (
	"bytes"
	"face_management/logger"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogRequest() gin.HandlerFunc {

	return func(c *gin.Context) {

		var requestBodyString string

		startTime := time.Now()

		requestMethod := c.Request.Method

		if requestMethod == "POST" || requestMethod == "PATCH" {

			requestBody, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				logger.Log.Error("GetRawData returned error!", []zapcore.Field{zap.String("Error:", err.Error())}...)
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
				requestBodyString = string(requestBody)
			}

		}

		c.Next()

		apiLatency := int64(time.Since(startTime) / time.Microsecond)

		if requestMethod == "POST" || requestMethod == "PATCH" {

			fields := []zapcore.Field{
				zap.String("x-request-id", c.Request.Header.Get("x-request-id")),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", requestMethod),
				zap.String("body", requestBodyString),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Int64("latency", apiLatency),
			}

			logger.Log.Info("api_stats", fields...)

		} else {

			fields := []zapcore.Field{
				zap.String("x-request-id", c.Request.Header.Get("x-request-id")),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", requestMethod),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Int64("latency", apiLatency),
			}

			logger.Log.Info("api_stats", fields...)

		}

	}

}
