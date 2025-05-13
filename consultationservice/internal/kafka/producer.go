package kafka

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/utils"
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	converter *utils.Converter
	conn      *kafka.Conn
}

func NewProducer(env *bootstrap.Env) (*Producer, error) {
	p := &Producer{}
	err := p.InitKafkaProducer(env)
	return p, err
}
func (r *Producer) InitKafkaProducer(env *bootstrap.Env) error {
	address := env.KafkaAddr
	topic := env.KafkaTopic
	partitionStr := env.KafkaPart
	if address == "" || topic == "" || partitionStr == "" {
		return errors.New("missing environment variables for kafka")
	}
	partition, err := r.converter.StringToInt(partitionStr)
	if err != nil {
		return err
	}
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		return errors.New("failed to connect to kafka leader: " + err.Error())
	}
	r.conn = conn
	log.Println("Kafka producer initialized successfully")
	return nil
}

func (r *Producer) ProducerSendMessage(message string) error {
	if r.conn == nil {
		return errors.New("kafka connection is not initialized")
	}
	r.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := r.conn.WriteMessages(kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Println("Failed to send kafka message:", err)
		return err
	}
	log.Println("Message sent to Kafka successfully")
	return nil
}
