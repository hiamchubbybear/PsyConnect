package kafka

import (
	"consultationservice/bootstrap"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	brokers            []string
	defaultTopic       string
	writer             *kafka.Writer
	loggingWriter      *kafka.Writer
	notificationWriter *kafka.Writer
}

func NewProducer(env *bootstrap.Env) (*Producer, error) {
	if env.KafkaAddr == "" || env.KafkaTopic == "" {
		log.Println("⚠️ Warning: Kafka configuration missing, running without Kafka producer")
		return &Producer{}, nil
	}

	brokers := []string{env.KafkaAddr}

	// Default consultation.test
	defaultWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    env.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})
	notificationWriter := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  brokers,
			Topic:    env.NotificationTopic,
			Balancer: &kafka.LeastBytes{},
		},
	)
	loggingWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "logging-service",
		Balancer: &kafka.LeastBytes{},
	})

	log.Println("Kafka producers initialized with topics:", env.KafkaTopic, "and logging-service")

	return &Producer{
		brokers:            brokers,
		writer:             defaultWriter,
		defaultTopic:       env.KafkaTopic,
		loggingWriter:      loggingWriter,
		notificationWriter: notificationWriter,
	}, nil
}

func (p *Producer) SendMessage(message string) error {
	if p.writer == nil {
		log.Println("⚠️ Kafka producer not initialized, skipping message")
		return nil
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
		log.Println("⚠️ Kafka notification producer not initialized, skipping message")
		return nil
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

func (p *Producer) SendToTopic(topic string, message string) error {
	if p.writer == nil {
		log.Println("⚠️ Kafka producer not initialized, skipping message")
		return nil
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  p.brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send Kafka message to specific topic %s: %v", topic, err)
		return err
	}

	log.Printf("Message sent to specific Kafka topic '%s' successfully", topic)
	return nil
}

func (p *Producer) SendLogs(message string) error {
	if p.loggingWriter == nil {
		log.Println("⚠️ Kafka logging producer not initialized, skipping message")
		return nil
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

func (p *Producer) SendSessionEvent(eventType string, payload interface{}) error {
	event := struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}{
		Key: "session_event",
		Value: struct {
			Type    string      `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type:    eventType,
			Payload: payload,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Failed to marshal session event: %v", err)
		return err
	}

	return p.SendNotification(string(data))
}

func (p *Producer) SendIncomingCallEvent(payload interface{}) error {
	return p.SendSessionEvent("consultation.incoming_call", payload)
}
