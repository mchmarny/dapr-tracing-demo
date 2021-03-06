package main

import (
	"encoding/json"
	"net/http"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
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

func rootHandler(c *gin.Context) {

	proto := c.GetHeader("x-forwarded-proto")
	if proto == "" {
		proto = "http"
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"host":    c.Request.Host,
		"proto":   proto,
		"version": AppVersion,
	})

}

func topicListHandler(c *gin.Context) {
	topics := []string{subTopic}
	logger.Printf("subscription topics: %v", topics)
	c.JSON(http.StatusOK, topics)
}

func subscribeHandler(c *gin.Context) {
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

	e := ce.NewEvent()
	if err := c.ShouldBindJSON(&e); err != nil {
		logger.Printf("error binding event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error processing your request, see logs for details",
		})
		return
	}

	var m SimpleMessage
	if err := json.Unmarshal(e.Data(), &m); err != nil {
		logger.Printf("error parsing event data (%s): %v", string(e.Data()), err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}
	logger.Printf("message: %v", m)

	// send
	if _, err := daprClient.InvokeBinding(ctx, bindingName, m); err != nil {
		logger.Printf("error binding output message %s (%v): %v", bindingName, m, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}
	logger.Print("sent to output binding")

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
