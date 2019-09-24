package main

import (
	"context"
	"flag"
	"fmt"
	"kafkaAPI/kafkaUtils"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

var (
	// kafka
	kafkaBrokerURL     string
	kafkaVerbose       bool
	kafkaTopicIn       string
	kafkaTopicOut      string
	kafkaConsumerGroup string
	kafkaClientID      string
)

var (
// kafka
)

func main() {
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaTopicIn, "kafka-topicIn", "foo", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaTopicOut, "kafka-topicOut", "foo2", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientID, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.Parse()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(kafkaBrokerURL, ",")

	// make a new reader that consumes from topic-A
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientID,
		Topic:           kafkaTopicIn,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)

	// connect to kafka
	kafkaProducer, err := kafkaUtils.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopicOut)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("unable to configure kafkaProducer")
		return
	}

	defer kafkaProducer.Close()
	defer reader.Close()

	for {
		log.Debug().Msg("Inside for loop1")
		m, err := reader.ReadMessage(context.Background())

		log.Debug().Msg("Inside for loop2")
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
			continue
		}

		log.Debug().Msg("Inside for loop2")
		value := m.Value
		log.Debug().Msg("Got a Message")
		//		if m.CompressionCodec == snappy.NewCompressionCodec() {
		//			_, err = snappy.NewCompressionCodec().Decode(value, m.Value)
		//		}

		var ctx = context.Background()
		err = kafkaUtils.Push(ctx, nil, m.Value)
		if err != nil {
			log.Error().Msg("Kafka write to topic Out failed")
		}

		log.Debug().Msg("Inside for loop3")
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
			continue
		}
		log.Debug().Msg("Printing Msg4")
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(value))
	}

	log.Debug().Msgf("Closing Down")

}
