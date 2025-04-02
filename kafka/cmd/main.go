package main

import (
	config "kafka/config"
	producer "kafka/producer"
	"log"
)

func main() {
	config.LoadConfig()
	err := producer.InitKafkaProducer()
	if err != nil {
		log.Fatal(err)
	}
	defer producer.CloseKafkaProducer()
}
