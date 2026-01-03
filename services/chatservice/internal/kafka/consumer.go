package kafka

import (
	"chatservice/bootstrap"
	"chatservice/internal/chat/model/model"
	"chatservice/internal/chat/repository/repository"
	"chatservice/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	converter  *utils.Converter
	conn       *kafka.Conn
	reader     *kafka.Reader
	repository repository.ChatRepository
}

func NewConsumer(env *bootstrap.Env) (*Consumer, error) {
	c := &Consumer{
		converter: &utils.Converter{},
	}
	err := c.initKafkaConsumer(env)
	return c, err
}

func (c *Consumer) initKafkaConsumer(env *bootstrap.Env) error {
	address := env.KafkaAddr
	topic := env.KafkaTopic
	partitionStr := env.KafkaPart

	if address == "" || topic == "" || partitionStr == "" {
		return errors.New("missing Kafka environment variables")
	}

	partition, err := c.converter.StringToInt(partitionStr)
	if err != nil {
		return err
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka leader: %v", err)
	}

	c.conn = conn
	log.Println("Kafka consumer initialized successfully")
	return nil
}

func (c *Consumer) ConsumeMessages() error {
	if c.conn == nil {
		return errors.New("Kafka connection is not initialized")
	}

	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.conn.ReadBatch(10e3, 1e6)
	defer func() {
		batch.Close()
		c.conn.Close()
	}()

	buffer := make([]byte, 10e3)

	for {
		n, err := batch.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(string(buffer[:n]))
	}
	return nil
}
func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("[Kafka Consumer] read error:", err)
			continue
		}

		var chat model.Chat
		if err := json.Unmarshal(m.Value, &chat); err != nil {
			log.Println("[Kafka Consumer] unmarshal error:", err)
			continue
		}

		_, err = c.repository.CreateChat(&chat)
		if err != nil {
			log.Println("[Kafka Consumer] failed to save chat:", err)
		}
	}
}
