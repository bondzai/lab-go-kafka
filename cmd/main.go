package main

import (
	"fmt"
	"os"
	"time"

	"github.com/introbond/lab-go-kafka/consumer"
	"github.com/introbond/lab-go-kafka/producer"
)

func main() {
	topic := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		err := consumer.ConsumeMessages(topic)
		if err != nil {
			fmt.Printf("Consumer error: %s\n", err)
		}
	}()

	go func() {
		for range ticker.C {
			err := fetchAndProduceBitcoinPrice(topic)
			if err != nil {
				fmt.Printf("Producer error: %s\n", err)
			}
		}
	}()

	// Prevent the main function from exiting immediately
	select {}
}

func fetchAndProduceBitcoinPrice(topic string) error {
	priceUSD, err := producer.FetchBitcoinPrice()
	if err != nil {
		return err
	}

	err = producer.ProduceMessage(topic, fmt.Sprintf("Bitcoin price in USD: %.2f", priceUSD))
	if err != nil {
		return err
	}

	fmt.Printf("Bitcoin price fetched and produced to Kafka: %.2f USD\n", priceUSD)
	return nil
}
