package main

import (
	"github.com/IBM/sarama"
	"kafka/config"
	"log"
)

var (
	kafkaProducer sarama.SyncProducer
	configKafka   = config.LoadConfig("../config")
)

func InitConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := configKafka.Kafka.Broker
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	for _, topic := range configKafka.Kafka.Topics {
		consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			log.Panicf("Error consuming topic %s: %v", topic, err)
		}
		go func(topic string, consumer sarama.PartitionConsumer) {
			for msg := range consumer.Messages() {
				log.Printf("Received message from %s: %s", topic, string(msg.Value))
			}
		}(topic, consumer)
	}
	select {}
}
func main() {
	InitConsumer()
}
