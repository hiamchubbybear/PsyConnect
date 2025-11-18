package kafka

import (
	"consultationservice/bootstrap"
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	defaultTopic       string
	writer             *kafka.Writer
	loggingWriter      *kafka.Writer
	notificationWriter *kafka.Writer
}

func NewProducer(env *bootstrap.Env) (*Producer, error) {
	if env.KafkaAddr == "" || env.KafkaTopic == "" {
		return nil, errors.New("missing Kafka configuration in environment")
	}

	// Default consultation.test
	defaultWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{env.KafkaAddr},
		Topic:    env.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})
	notificationWriter := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  []string{env.KafkaAddr},
			Topic:    env.NotificationTopic,
			Balancer: &kafka.LeastBytes{},
		},
	)
	loggingWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{env.KafkaAddr},
		Topic:    "logging-service",
		Balancer: &kafka.LeastBytes{},
	})

	log.Println("Kafka producers initialized with topics:", env.KafkaTopic, "and logging-service")

	return &Producer{
		writer:             defaultWriter,
		defaultTopic:       env.KafkaTopic,
		loggingWriter:      loggingWriter,
		notificationWriter: notificationWriter,
	}, nil
}

func (p *Producer) SendMessage(message string) error {
	if p.writer == nil {
		return errors.New("default Kafka writer is not initialized")
	}

	err := p.writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send Kafka message to topic %s: %v", p.defaultTopic, err)
		return err
	}

	log.Printf("Message sent to Kafka topic '%s' successfully", p.defaultTopic)
	return nil
}
func (p *Producer) SendNotification(message string) error {
	if p.notificationWriter == nil {
		return errors.New("default Kafka writer is not initialized")
	}

	err := p.notificationWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send Kafka message to topic %s: %v", p.notificationWriter.Topic, err)
		return err
	}

	log.Printf("Message sent to Kafka topic '%s' successfully", p.notificationWriter.Topic)
	return nil
}
func (p *Producer) SendLogs(message string) error {
	if p.loggingWriter == nil {
		return errors.New("logging Kafka writer is not initialized")
	}

	err := p.loggingWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send Kafka message to topic logging-service: %v", err)
		return err
	}

	log.Printf("Message sent to Kafka topic 'logging-service' successfully")
	return nil
}

func (p *Producer) Close() error {
	var err1, err2 error
	if p.writer != nil {
		err1 = p.writer.Close()
	}
	if p.loggingWriter != nil {
		err2 = p.loggingWriter.Close()
	}

	if err1 != nil {
		return err1
	}
	return err2
}
