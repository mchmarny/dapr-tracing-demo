package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestBindHandler(t *testing.T) {

	daprClient = GetTestClient()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/", receiveHandler)
	w := httptest.NewRecorder()

	in := &SimpleMessage{
		ID:        uuid.New().String(),
		Text:      "test",
		CreatedOn: time.Now(),
	}
	data, _ := json.Marshal(in)

	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func GetTestClient() *TestClient {
	return &TestClient{}
}

var (
	_ = Client(&TestClient{})
)

// TestClient is a test dapr client
type TestClient struct {
}

// SaveState mocks saving state
func (c *TestClient) SaveState(ctx trace.SpanContext, store, key string, data interface{}) error {
	return nil
}

// InvokeService mocks invoking service
func (c *TestClient) InvokeService(ctx trace.SpanContext, service, method string, data interface{}) (out []byte, err error) {
	in := &SimpleMessage{
		ID:        uuid.New().String(),
		Text:      "test",
		CreatedOn: time.Now(),
	}
	b, _ := json.Marshal(in)
	return b, nil
}

// Publish mocks publishing to topic
func (c *TestClient) Publish(ctx trace.SpanContext, topic string, data interface{}) error {
	return nil
}
