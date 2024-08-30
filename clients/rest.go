package clients

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type restClient struct{}

type restClientInterface interface {
	Get(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error)
	Post(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error)
	Patch(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error)
	Delete(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error)
}

var RestClient restClientInterface

func init() {
	RestClient = &restClient{}
}

func (rc *restClient) Get(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	return statusCode, responseBody, nil

}

func (rc *restClient) Post(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	return statusCode, responseBody, nil

}

func (rc *restClient) Patch(url string, headers map[string]string, body interface{}, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	return statusCode, responseBody, nil

}

func (rc *restClient) Delete(url string, headers map[string]string, timeout time.Duration) (statusCode int, responseBody []byte, err error) {

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	return statusCode, responseBody, nil

}
