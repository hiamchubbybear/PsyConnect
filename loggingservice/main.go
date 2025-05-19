package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func listenOnKafkaHandler() error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9092",
		"group.id":          "logging-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	topics := []string{
		"logging-service",
		"identity-service",
		"profile-service",
		"friends-service",
	}
	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening to topics:", topics)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
			run = false

		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				topic := *e.TopicPartition.Topic
				fmt.Printf("\nMessage from topic [%s], partition [%d], offset [%d]\n",
					topic, e.TopicPartition.Partition, e.TopicPartition.Offset)

				if len(e.Key) > 0 {
					fmt.Printf("Key: %s\n", string(e.Key))
				}

				var jsonData map[string]interface{}
				if err := json.Unmarshal(e.Value, &jsonData); err == nil {
					jsonString, _ := json.MarshalIndent(jsonData, "", "  ")
					fmt.Printf("JSON Content:\n%s\n", string(jsonString))
					_ = writeToLog(string(jsonString))
				} else {
					fmt.Printf("Raw Content: %s\n", string(e.Value))
					_ = writeToLog(string(e.Value))
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Kafka Error: %v\n", e)
			}
		}
	}
	return nil
}
func main() {
	listenOnKafkaHandler()
}
func writeToLog(message string) error {
	// Parse message thành map
	var event map[string]interface{}
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		return err
	}

	// Tạo logger như cũ
	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    15,
		MaxBackups: 10,
		MaxAge:     120,
		Compress:   true,
		LocalTime:  true,
	})
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "timestamp",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		logWriter,
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()

	get := func(k string) string {
		if val, ok := event[k]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
		return ""
	}
	log.Printf("Log level %s", get("level"))
	if get("level") == "LOG" || get("level") == "AUDIT" {
		logger.Info("LogEvent",
			zap.String("service", get("service")),
			zap.String("action", get("action")),
			zap.String("userId", get("userId")),
			zap.String("traceId", get("traceId")),
			zap.String("timestamp", get("timestamp")),
			zap.Any("metadata", event["metadata"]),
			zap.String("message", get("message")))
	} else if get("level") == "ERROR" || get("level") == "WARN" {
		logger.Error("LogEvent",
			zap.String("service", get("service")),
			zap.String("action", get("action")),
			zap.String("userId", get("userId")),
			zap.String("traceId", get("traceId")),
			zap.String("timestamp", get("timestamp")),
			zap.Any("metadata", event["metadata"]),
			zap.String("message", get("message")))
	}

	return nil
}
