package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func decorateHandler(c *gin.Context) {
	httpFmt := tracecontext.HTTPFormat{}
	ctx, ok := httpFmt.SpanContextFromRequest(c.Request)
	if !ok {
		ctx = trace.SpanContext{}
	}

	logger.Printf("Trace Info: 0-%x-%x-%x",
		ctx.TraceID[:],
		ctx.SpanID[:],
		[]byte{byte(ctx.TraceOptions)})

	m := SimpleMessage{}
	if err := c.ShouldBindJSON(&m); err != nil || m.Text == "" {
		logger.Printf("error binding request: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	m.Text = fmt.Sprintf("**%s** -- decorated", m.Text)

	c.JSON(http.StatusOK, m)
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
