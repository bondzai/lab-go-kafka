package main

import (
	"fmt"
	"os"

	"github.com/introbond/lab-go-kafka/consumer"
	"github.com/introbond/lab-go-kafka/producer"
)

func main() {
	topic := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")
	message := "hello from local"

	// Start Producer Goroutine
	go func() {
		err := producer.ProduceMessage(topic, message)
		if err != nil {
			fmt.Printf("Producer error: %s\n", err)
		}
	}()

	// Start Consumer Goroutine
	go func() {
		err := consumer.ConsumeMessages(topic)
		if err != nil {
			fmt.Printf("Consumer error: %s\n", err)
		}
	}()

	// Prevent the main function from exiting immediately
	select {}
}
