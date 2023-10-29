package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Retrieve Kafka configuration from environment variables
	bootstrapServers := os.Getenv("CLOUDKARAFKA_BROKERS")
	username := os.Getenv("CLOUDKARAFKA_USERNAME")
	password := os.Getenv("CLOUDKARAFKA_PASSWORD")
	topic := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")

	// Producer
	producerConfig := kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     username,
		"sasl.password":     password,
	}

	p, err := kafka.NewProducer(&producerConfig)
	if err != nil {
		fmt.Printf("Error creating producer: %s\n", err)
		return
	}

	defer p.Close()

	// Produce a "hello" message
	message := "hello from local"

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		fmt.Printf("Error producing message: %s\n", err)
	} else {
		fmt.Println("Message sent:", message)
	}

	p.Flush(15 * 1000) // Wait for any outstanding messages to be delivered and delivery reports to be received.

	// Consumer
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     username,
		"sasl.password":     password,
		"group.id":          "my-group",
		"auto.offset.reset": "earliest", // Start consuming from the beginning of the topic.
	}

	c, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		fmt.Printf("Error creating consumer: %s\n", err)
		return
	}

	defer c.Close()

	// Subscribe to the topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %s\n", err)
		return
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error reading message: %v\n", err)
		}
	}
}
