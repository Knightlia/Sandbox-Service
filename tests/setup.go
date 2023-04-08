package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sandbox-service/app"
	"sandbox-service/cache"
)

var (
	testConn  *websocket.Conn
	testToken string
)

func SetupTests() *httptest.Server {
	cache.InitCaches()

	a := app.NewApp()
	a.InitApp()
	a.InitRoutes()

	s := httptest.NewServer(a.Chi)

	// A test user for mock data
	testConn, testToken = connectTestClient(s)
	cache.SessionCache.Set(testConn, testToken, ttlcache.DefaultTTL)
	cache.UserCache.Set(testToken, "existing-nickname", ttlcache.DefaultTTL)

	return s
}

func TeardownTests(body io.ReadCloser) {
	_ = body.Close()
	_ = testConn.Close(websocket.StatusNormalClosure, "")
}

func GET(s *httptest.Server, path string) *http.Response {
	return doRequest(s.URL+path, http.MethodGet, nil, nil)
}

func POST(s *httptest.Server, path string, headers map[string]string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	return doRequest(s.URL+path, http.MethodPost, headers, bytes.NewBuffer(b))
}

func DefaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"token":        testToken,
	}
}

func AssertJSONResponse(t *testing.T, expected interface{}, actual io.ReadCloser) {
	expectedBytes, _ := json.Marshal(expected)
	actualBytes, _ := io.ReadAll(actual)
	assert.JSONEq(t, string(expectedBytes), string(actualBytes))
}

// ---

func doRequest(url string, method string, headers map[string]string, body io.Reader) *http.Response {
	req, _ := http.NewRequest(method, url, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, _ := http.DefaultClient.Do(req)
	return res
}

func connectTestClient(s *httptest.Server) (*websocket.Conn, string) {
	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	c, _, _ := websocket.Dial(context.Background(), url, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	payload := make(map[string]string, 0)
	_ = wsjson.Read(ctx, c, &payload)

	return c, payload["t"]
}
