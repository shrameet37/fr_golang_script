package test

import (
	"encoding/json"
	"face_management/app"
	"face_management/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHealthNoErrors(t *testing.T) {

	assert := assert.New(t)

	app.SetupHealthRoute()

	resWriter := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/health", nil)

	req.Header.Set("Content-Type", "application/json")

	app.Router.ServeHTTP(resWriter, req)

	healthResponse := models.SuccessResponse{}
	json.Unmarshal(resWriter.Body.Bytes(), &healthResponse)

	assert.Equal(200, resWriter.Code)
	assert.Equal("success", healthResponse.Type)
	assert.Equal("healthy", healthResponse.Message)

}
