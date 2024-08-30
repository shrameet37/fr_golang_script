package test

import (
	"face_management/clients"
	"time"
)

var (
	funcRestClientMockResponse func() (statusCode int, responseBody []byte, err error)
)

type restClientMock struct{}

func init() {

	clients.RestClient = &restClientMock{}

}

func (rcm *restClientMock) Get(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	return funcRestClientMockResponse()

}

func (rcm *restClientMock) Post(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	return funcRestClientMockResponse()

}

func (rcm *restClientMock) Patch(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	return funcRestClientMockResponse()

}

func (rcm *restClientMock) Delete(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	return funcRestClientMockResponse()

}
