package main

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

func createProducer(config kafka.ConfigMap) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %s", err)
	}
	return producer, nil
}

func createConsumer(config kafka.ConfigMap, groupID string) (*kafka.Consumer, error) {
	config["group.id"] = groupID
	config["auto.offset.reset"] = "earliest"
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer: %s", err)
	}
	return consumer, nil
}

func sendMessage(producer *kafka.Producer, topic, message string) error {
	err := producer.Produce(&kafka.Message{
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

func consumeMessages(consumer *kafka.Consumer, topics []string) error {
	err := consumer.SubscribeTopics(topics, nil)
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

func main() {
	config := getKafkaConfig()
	topic := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")

	producer, err := createProducer(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()

	message := "hello from local"
	err = sendMessage(producer, topic, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	consumer, err := createConsumer(config, topic)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer consumer.Close()

	err = consumeMessages(consumer, []string{topic})
	if err != nil {
		fmt.Println(err)
		return
	}
}
