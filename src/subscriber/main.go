package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	dapr "github.com/mchmarny/godapr/v1"
	"go.opencensus.io/trace"
	"gopkg.in/olahol/melody.v1"
)

var (
	logger = log.New(os.Stdout, "SUBSCRIBER == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8083")
	subTopic    = env.MustGetEnvVar("SUBSCRIBER_TOPIC_NAME", "messages")
	bindingName = env.MustGetEnvVar("PRODUCER_BINDING_NAME", "send")

	broadcaster *melody.Melody

	// dapr
	daprClient Client

	// test client against local interace
	_ = Client(dapr.NewClient())
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	daprClient = dapr.NewClient()

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// simple routes
	r.GET("/", subscribeHandler)
	r.GET("/dapr/subscribe", topicListHandler)

	// topic route
	subRoute := fmt.Sprintf("/%s", subTopic)
	logger.Printf("subscription route: %s", subRoute)
	r.POST(subRoute, subscribeHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := http.ListenAndServe(hostPort, r); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

// Options midleware
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Client is the minim client support for testing
type Client interface {
	InvokeBinding(ctx trace.SpanContext, binding string, in interface{}) (out []byte, err error)
}
