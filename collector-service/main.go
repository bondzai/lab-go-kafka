package main

import (
	"fmt"
	"time"

	"github.com/introbond/lab-go-kafka/constants"
	"github.com/introbond/lab-go-kafka/consumer"
	"github.com/introbond/lab-go-kafka/producer"
)

func main() {

	producerTicker := time.NewTicker(3 * time.Second)

	go func() {
		err := consumer.ConsumeMessages(constants.KafkaCryptoPriceTopic)
		if err != nil {
			fmt.Printf("Consumer error: %s\n", err)
		}
	}()

	go func() {
		for range producerTicker.C {
			err := producer.ProduceCryptoPriceMessages(constants.KafkaCryptoPriceTopic)
			if err != nil {
				fmt.Printf("Producer error: %s\n", err)
			}
		}
	}()

	select {}
}
