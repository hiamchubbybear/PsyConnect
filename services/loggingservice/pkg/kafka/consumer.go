package kafka

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/loggingservice/pkg/models"
	"github.com/loggingservice/pkg/settings"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)


type Consumer struct {
	reader   *kafka.Reader
	config   settings.KafkaConfig
	logger   *zap.Logger
	handlers []LogHandler
}


type LogHandler func(*models.LogEvent) error


func NewConsumer(config settings.KafkaConfig, logger *zap.Logger) (*Consumer, error) {
	brokers := strings.Split(config.BootstrapServers, ",")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        config.GroupID,
		GroupTopics:    config.Topics,
		MinBytes:       10e3, 
		MaxBytes:       10e6, 
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.FirstOffset,
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			logger.Error(fmt.Sprintf(msg, args...))
		}),
	})

	return &Consumer{
		reader:   reader,
		config:   config,
		logger:   logger,
		handlers: make([]LogHandler, 0),
	}, nil
}


func (c *Consumer) AddHandler(handler LogHandler) {
	c.handlers = append(c.handlers, handler)
}


func (c *Consumer) Start() error {
	c.logger.Info("Kafka consumer started",
		zap.Strings("topics", c.config.Topics),
		zap.String("group_id", c.config.GroupID),
	)

	
	ctx := context.Background()

	for {
		
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				c.logger.Info("Reader closed")
				return nil
			}
			c.logger.Error("Failed to read message", zap.Error(err))
			
			time.Sleep(1 * time.Second)
			continue
		}

		c.handleMessage(m)
	}
}


func (c *Consumer) handleMessage(msg kafka.Message) {
	topic := msg.Topic

	c.logger.Debug("Received message",
		zap.String("topic", topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
	)

	
	event, err := models.ParseLogEvent(msg.Value)
	if err != nil {
		c.logger.Error("Failed to parse log event",
			zap.Error(err),
			zap.String("topic", topic),
			zap.ByteString("raw_message", msg.Value),
		)
		return
	}

	
	if event.Service == "" {
		event.Service = topic
	}

	
	for _, handler := range c.handlers {
		if err := handler(event); err != nil {
			c.logger.Error("Handler failed",
				zap.Error(err),
				zap.String("service", event.Service),
				zap.String("level", event.Level),
			)
		}
	}
}


func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka consumer")
	return c.reader.Close()
}
