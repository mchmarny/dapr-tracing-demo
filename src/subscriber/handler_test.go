package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/", subscribeHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./event.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTopicListHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/", topicListHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var out []string
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Len(t, out, 1)

}
