package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/loggingservice/pkg/models"
	"github.com/loggingservice/pkg/settings"
	"go.uber.org/zap"
)

// Consumer represents a Kafka consumer for log events
type Consumer struct {
	consumer *kafka.Consumer
	config   settings.KafkaConfig
	logger   *zap.Logger
	handlers []LogHandler
}

// LogHandler is a function that processes log events
type LogHandler func(*models.LogEvent) error

// NewConsumer creates a new Kafka consumer
func NewConsumer(config settings.KafkaConfig, logger *zap.Logger) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupID,
		"auto.offset.reset": config.AutoOffsetReset,
		// Performance tuning
		"session.timeout.ms":          6000,
		"max.poll.interval.ms":        300000,
		"enable.auto.commit":          true,
		"auto.commit.interval.ms":     5000,
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: c,
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
	// Subscribe to topics
	err := c.consumer.SubscribeTopics(c.config.Topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	c.logger.Info("Kafka consumer started",
		zap.Strings("topics", c.config.Topics),
		zap.String("group_id", c.config.GroupID),
	)

	// Setup signal handling for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			c.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
			run = false

		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				c.handleMessage(e)

			case kafka.Error:
				c.logger.Error("Kafka error", zap.Error(e))
				// Check if it's a fatal error
				if e.Code() == kafka.ErrAllBrokersDown {
					c.logger.Fatal("All Kafka brokers are down")
					run = false
				}

			default:
				// Ignore other events
			}
		}
	}

	return nil
}

// handleMessage processes a Kafka message
func (c *Consumer) handleMessage(msg *kafka.Message) {
	topic := *msg.TopicPartition.Topic

	c.logger.Debug("Received message",
		zap.String("topic", topic),
		zap.Int32("partition", msg.TopicPartition.Partition),
		zap.Int64("offset", int64(msg.TopicPartition.Offset)),
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
	return c.consumer.Close()
}
