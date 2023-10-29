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

	producerTicker := time.NewTicker(300 * time.Second)

	go func() {
		err := consumer.ConsumeMessages(topic)
		if err != nil {
			fmt.Printf("Consumer error: %s\n", err)
		}
	}()

	go func() {
		for range producerTicker.C {
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

	message := fmt.Sprintf("%d - Bitcoin price in USD: %.2f", time.Now().Unix(), priceUSD)

	err = producer.ProduceMessage(topic, message)
	if err != nil {
		return err
	}

	fmt.Printf("Bitcoin price fetched and produced to Kafka: %.2f USD\n", priceUSD)
	return nil
}
