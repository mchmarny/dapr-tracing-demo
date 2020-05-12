package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/dapr-tracing-demo/producer/message"
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

func caseHandler(c *gin.Context) {
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

	m := message.SimpleMessage{}
	if err := c.ShouldBindJSON(&m); err != nil || m.Text == "" {
		logger.Printf("error binding request: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	m.Text = fmt.Sprintf("%s -- decorated", strings.ToUpper(m.Text))

	c.JSON(http.StatusOK, m)
}
