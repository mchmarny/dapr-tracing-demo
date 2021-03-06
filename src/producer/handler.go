package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
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

func receiveHandler(c *gin.Context) {
	logger.Printf("traceparent: %s", c.GetHeader("traceparent"))
	logger.Printf("tracestate: %s", c.GetHeader("tracestate"))

	httpFmt := tracecontext.HTTPFormat{}
	ctx, ok := httpFmt.SpanContextFromRequest(c.Request)
	if !ok {
		ctx = trace.SpanContext{}
	}

	logger.Printf("trace info: 0-%x-%x-%x",
		ctx.TraceID[:],
		ctx.SpanID[:],
		[]byte{byte(ctx.TraceOptions)})

	var m SimpleMessage
	if err := c.ShouldBindJSON(&m); err != nil {
		logger.Printf("error binding input message: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("raw: %v", m)

	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.CreatedOn.IsZero() {
		m.CreatedOn = time.Now()
	}

	// save
	err := daprClient.SaveState(ctx, stateStore, m.ID, m)
	if err != nil {
		logger.Printf("error saving state: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}
	logger.Printf("saved: %v", m)

	// format
	b, err := daprClient.InvokeService(ctx, serviceName, serviceMethod, m)
	if err != nil {
		logger.Printf("error invoking service (%s/%s): %v",
			serviceName, serviceMethod, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	var d SimpleMessage
	if err := json.Unmarshal(b, &d); err != nil {
		logger.Printf("error parsing service response (%s): %v", string(b), err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}
	logger.Printf("formatted: %v", m)

	// publish
	if err = daprClient.Publish(ctx, eventTopic, d); err != nil {
		logger.Printf("error publishing content %s (%v): %v", eventTopic, d, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}
	logger.Print("published")

	c.JSON(http.StatusOK, gin.H{})
}

// SimpleMessage represents simple message
type SimpleMessage struct {
	// ID is the ID of the message
	ID string `json:"id,omitempty"`
	// Text is the test of the message
	Text string `json:"txt,omitempty"`
	// CreatedOn is the time when this message was created
	CreatedOn time.Time `json:"on,omitempty"`
}

// Envelope is the output binding content
type Envelope []struct {
	Data interface{} `json:"data,omitempty"`
}
