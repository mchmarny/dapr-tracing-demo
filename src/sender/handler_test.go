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
)

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/", postHandler)
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
