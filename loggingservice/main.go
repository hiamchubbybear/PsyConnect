package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
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
				} else {
					fmt.Printf("Raw Content: %s\n", string(e.Value))
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Kafka Error: %v\n", e)
			}
		}
	}
}
