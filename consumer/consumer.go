package consumer

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func getKafkaConfig() kafka.ConfigMap {
	bootstrapServers := os.Getenv("CLOUDKARAFKA_BROKERS")
	username := os.Getenv("CLOUDKARAFKA_USERNAME")
	password := os.Getenv("CLOUDKARAFKA_PASSWORD")

	return kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     username,
		"sasl.password":     password,
	}
}

func ConsumeMessages(topic string) error {
	config := getKafkaConfig()
	config["group.id"] = topic
	config["auto.offset.reset"] = "earliest"

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return fmt.Errorf("error creating consumer: %s", err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return fmt.Errorf("error subscribing to topics: %s", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error reading message: %v\n", err)
		}
	}
}
