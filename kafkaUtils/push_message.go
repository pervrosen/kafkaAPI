package kafkaUtils

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

//Push {value} onto {key}
func Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(parent, message)
}
