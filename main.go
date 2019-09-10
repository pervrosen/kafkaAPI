package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"kafkaAPI/kafkaUtils"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	//"github.com/segmentio/kafka-go/snappy"
)

var logger = log.With().Str("pkg", "main").Logger()

var (
	listenAddrAPI string

	// kafka
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaClientID  string
	kafkaTopic     string
)

func main() {
	flag.StringVar(&listenAddrAPI, "listen-address", "0.0.0.0:9000", "Listen address for api")
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaClientID, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic to push")
	flag.Parse()

	// connect to kafka
	kafkaProducer, err := kafkaUtils.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("unable to configure kafka")
		return
	}
	defer kafkaProducer.Close()

	var errChan = make(chan error, 1)

	go func() {
		log.Info().Msgf("starting server at %s", listenAddrAPI)
		errChan <- server(listenAddrAPI)
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.Info().Msg("got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			logger.Error().Err(err).Msg("error while running api, exiting...")
		}
	}
}

func server(listenAddr string) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.POST("/api/v1/data", postDataToKafka)

	router.POST("/api/v1/test", testar)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	router.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")

		if user == "pvr" {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": "Superhero"})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// for debugging purpose
	for _, routeInfo := range router.Routes() {
		logger.Debug().
			Str("path", routeInfo.Path).
			Str("handler", routeInfo.Handler).
			Str("method", routeInfo.Method).
			Msg("registered routes")
	}

	return router.Run(listenAddr)
}

func testar(ctx *gin.Context) {
	user := "pvr"
	value := "lalala"
	ok := true
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	}
}

func postDataToKafka(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	form := &struct {
		Text string `form:"text" json:"text"`
	}{}

	ctx.Bind(form)
	formInBytes, err := json.Marshal(form)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while marshalling json: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	err = kafkaUtils.Push(parent, nil, formInBytes)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "success push data into kafka",
		"data":    form,
	})
}
