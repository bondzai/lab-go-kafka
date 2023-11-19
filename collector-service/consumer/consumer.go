package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/zenithero/lab-go-kafka/constants"
)

func ConsumeMessages(topic string) error {
	config := constants.GetKafkaConfig()
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
