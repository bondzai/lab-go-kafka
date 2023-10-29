package producer

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

func ProduceMessage(topic, message string) error {
	config := getKafkaConfig()

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return fmt.Errorf("error creating producer: %s", err)
	}
	defer producer.Close()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		return fmt.Errorf("error producing message: %s", err)
	}
	fmt.Println("Message sent:", message)
	producer.Flush(15 * 1000) // Wait for any outstanding messages to be delivered and delivery reports to be received.
	return nil
}
