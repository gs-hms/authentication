package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer is a wrapper around kafka-go's Writer to publish messages to Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new instance of a Kafka Producer.
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
	}
}

// Publish serializes the value to JSON and publishes it to the Kafka topic.
func (p *Producer) Publish(ctx context.Context, userID uint64, key string, value map[string]string) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message value: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: bytes,
		Headers: []kafka.Header{
			{Key: "user-id", Value: []byte(fmt.Sprintf("%d", userID))},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	log.Printf("published message to kafka | topic=%s | key=%s", p.writer.Topic, key)
	return nil
}

// Close closes the underlying Kafka writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
