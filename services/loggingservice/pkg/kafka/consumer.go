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

// Consumer represents a Kafka consumer for log events
type Consumer struct {
	reader   *kafka.Reader
	config   settings.KafkaConfig
	logger   *zap.Logger
	handlers []LogHandler
}

// LogHandler is a function that processes log events
type LogHandler func(*models.LogEvent) error

// NewConsumer creates a new Kafka consumer
func NewConsumer(config settings.KafkaConfig, logger *zap.Logger) (*Consumer, error) {
	brokers := strings.Split(config.BootstrapServers, ",")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        config.GroupID,
		GroupTopics:    config.Topics,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
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

// AddHandler adds a log handler
func (c *Consumer) AddHandler(handler LogHandler) {
	c.handlers = append(c.handlers, handler)
}

// Start starts consuming messages from Kafka
func (c *Consumer) Start() error {
	c.logger.Info("Kafka consumer started",
		zap.Strings("topics", c.config.Topics),
		zap.String("group_id", c.config.GroupID),
	)

	// Context for cancellation
	ctx := context.Background()

	for {
		// ReadMessage automatically commits offsets when using consumer groups
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				c.logger.Info("Reader closed")
				return nil
			}
			c.logger.Error("Failed to read message", zap.Error(err))
			// Exponential backoff could be added here
			time.Sleep(1 * time.Second)
			continue
		}

		c.handleMessage(m)
	}
}

// handleMessage processes a Kafka message
func (c *Consumer) handleMessage(msg kafka.Message) {
	topic := msg.Topic

	c.logger.Debug("Received message",
		zap.String("topic", topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
	)

	// Parse log event
	event, err := models.ParseLogEvent(msg.Value)
	if err != nil {
		c.logger.Error("Failed to parse log event",
			zap.Error(err),
			zap.String("topic", topic),
			zap.ByteString("raw_message", msg.Value),
		)
		return
	}

	// Ensure service name is set from topic if not present
	if event.Service == "" {
		event.Service = topic
	}

	// Call all handlers
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

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka consumer")
	return c.reader.Close()
}
