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

func bindHandler(c *gin.Context) {
	httpFmt := tracecontext.HTTPFormat{}
	ctx, ok := httpFmt.SpanContextFromRequest(c.Request)
	if !ok {
		ctx = trace.SpanContext{}
	}

	logger.Printf("Trace Info: 0-%x-%x-%x",
		ctx.TraceID[:],
		ctx.SpanID[:],
		[]byte{byte(ctx.TraceOptions)})

	var m SimpleMessage
	if err := c.ShouldBindJSON(&m); err != nil {
		logger.Printf("error binding message: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.CreatedOn.IsZero() {
		m.CreatedOn = time.Now()
	}

	// save original tweet in case we need to reprocess it
	err := daprClient.SaveState(ctx, stateStore, m.ID, m)
	if err != nil {
		logger.Printf("error saving state: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	// score simple tweet
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

	// publish simple tweet
	if err = daprClient.Publish(ctx, eventTopic, d); err != nil {
		logger.Printf("error publishing content (%v): %v", d, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

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
