package kafka

import (
	"consultationservice/bootstrap"
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	defaultTopic string
	writer       *kafka.Writer
}

func NewProducer(env *bootstrap.Env) (*Producer, error) {
	if env.KafkaAddr == "" || env.KafkaTopic == "" {
		return nil, errors.New("missing Kafka configuration in environment")
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{env.KafkaAddr},
		Topic:    env.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	log.Println("Kafka producer initialized with topic:", env.KafkaTopic)

	return &Producer{
		writer:       writer,
		defaultTopic: env.KafkaTopic,
	}, nil
}

func (p *Producer) SendMessage(message string) error {
	return p.sendToTopic(p.defaultTopic, message)
}

func (p *Producer) SendLogs(message string) error {
	return p.sendToTopic("logging-service", message)
}

func (p *Producer) sendToTopic(topic, message string) error {
	if p.writer == nil {
		return errors.New("Kafka writer is not initialized")
	}

	err := p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send Kafka message to topic %s: %v", topic, err)
		return err
	}

	log.Printf("Message sent to Kafka topic '%s' successfully", topic)
	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
