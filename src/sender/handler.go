package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      AppVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func postHandler(c *gin.Context) {
	logger.Printf("traceparent: %s", c.GetHeader("traceparent"))
	logger.Printf("tracestate: %s", c.GetHeader("tracestate"))

	m := SimpleMessage{}
	if err := c.ShouldBindJSON(&m); err != nil || m.Text == "" {
		logger.Printf("error binding request: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("message: %v", m)

	c.JSON(http.StatusOK, gin.H{})
}

// SimpleMessage represents simple message
type SimpleMessage struct {
	// ID is the ID of the message
	ID string `json:"id"`
	// Text is the test of the message
	Text string `json:"txt"`
	// CreatedOn is the time when this message was created
	CreatedOn time.Time `json:"on"`
}
