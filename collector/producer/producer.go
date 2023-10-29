package producer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/introbond/lab-go-kafka/constants"
)

func getKafkaConfig() kafka.ConfigMap {

	return kafka.ConfigMap{
		"bootstrap.servers": constants.KafkaBootstrapServers,
		"security.protocol": constants.KafkaSecureProtocol,
		"sasl.mechanism":    constants.KafkaSaslMech,
		"sasl.username":     constants.KafkaUsername,
		"sasl.password":     constants.KafkaPassword,
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

type CoinGeckoPriceResponse struct {
	Bitcoin struct {
		Usd float64 `json:"usd"`
	} `json:"bitcoin"`
}

func FetchBitcoinPrice() (float64, error) {
	coinGeckoURL := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
	resp, err := http.Get(coinGeckoURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	var priceResponse CoinGeckoPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
		return 0, err
	}

	return priceResponse.Bitcoin.Usd, nil
}

func ProduceBitcoinPriceMessage(topic string, priceUSD float64) error {
	message := fmt.Sprintf("Bitcoin price in USD: %.2f", priceUSD)
	return ProduceMessage(topic, message)
}
