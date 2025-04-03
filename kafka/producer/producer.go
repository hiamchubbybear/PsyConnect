package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"kafka/config"
	"log"
)

var (
	kafkaProducer sarama.SyncProducer
	configKafka   = config.LoadConfig("")
)

func InitKafkaProducer() error {
	broker := configKafka.Kafka.Broker
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewSyncProducer(broker, kafkaConfig)
	if err != nil {
		log.Fatalf("There some error while creates producer %v", err)
	}
	kafkaProducer = producer
	fmt.Printf("Kafka producer create success. ")
	return nil
}
func SendMessage(topic, key, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := kafkaProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Send message failed: %v", err)
	}
	fmt.Printf("Send message to  %s (Partition %d, Offset %d)\n", topic, partition, offset)
	return nil
}
func CloseKafkaProducer() {
	if kafkaProducer != nil {
		kafkaProducer.Close()
		fmt.Println("Kafka Producer closed")
	}
}
