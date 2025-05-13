package kafka

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Consumer struct {
	converter *utils.Converter
	conn      *kafka.Conn
}

func NewConsumer(env *bootstrap.Env) (*Producer, error) {
	p := &Producer{}
	err := p.InitKafkaProducer(env)
	return p, err
}
func (r *Consumer) InitKafkaConsumer(env *bootstrap.Env) error {
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
func (r *Consumer) ConsumerMessage() error {
	if r.conn == nil {
		return errors.New("kafka connection is not initialized")
	}
	r.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := r.conn.ReadBatch(10e3, 1e6)
	b := make([]byte, 10e3)
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := r.conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
	return nil
}
