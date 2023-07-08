package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knightlia/sandbox-service/app"
	"github.com/olahol/melody"
	"github.com/stretchr/testify/assert"
)

var validToken = "valid_token"

// GET performs a http GET request on the test server.
func GET(s *httptest.Server, path string) *http.Response {
	return doRequest(s.URL+path, http.MethodGet, nil, nil)
}

// POST performs a http POST request on the test server.
func POST(s *httptest.Server, path string, headers map[string]string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	return doRequest(s.URL+path, http.MethodPost, headers, bytes.NewBuffer(b))
}

// Sets up the test server.
func setup() (*httptest.Server, app.App) {
	a := app.NewApp(melody.New())
	a.InitApp()
	a.InitRoutes()

	a.UserCache.Store(validToken, "")

	return httptest.NewServer(a.Chi), a
}

// Cleans up resources from the tests.
func teardown(body io.ReadCloser) {
	_ = body.Close()
}

// Setup default headers with headers.
func defaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"token":        validToken,
	}
}

// Assert json response body.
func assertJSONResponse(t *testing.T, expected interface{}, actual io.ReadCloser) {
	expectedBytes, _ := json.Marshal(expected)
	actualBytes, _ := io.ReadAll(actual)
	assert.JSONEq(t, string(expectedBytes), string(actualBytes))
}

// Performs the actual http request and returns the response.
func doRequest(url string, method string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest(method, url, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, _ := http.DefaultClient.Do(req)
	return res
}
