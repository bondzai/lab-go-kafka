package producer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/introbond/lab-go-kafka/constants"
)

// Add more supported cryptocurrencies
var supportedCryptos = []string{"bitcoin", "ethereum", "the-graph"}

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
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
	TheGraph struct {
		Usd float64 `json:"usd"`
	} `json:"the-graph"`
}

func FetchCryptoPrice(cryptoName string) (float64, error) {
	coinGeckoURL := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", cryptoName)
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

	switch cryptoName {
	case "bitcoin":
		return priceResponse.Bitcoin.Usd, nil
	case "ethereum":
		return priceResponse.Ethereum.Usd, nil
	case "the-graph":
		return priceResponse.TheGraph.Usd, nil
	default:
		return 0, fmt.Errorf("unsupported cryptocurrency: %s", cryptoName)
	}
}

func ProduceCryptoPriceMessages(topic string) error {

	for _, crypto := range supportedCryptos {
		price, err := FetchCryptoPrice(crypto)
		if err != nil {
			fmt.Printf("Error fetching %s price: %s\n", crypto, err)
			continue
		}

		message := fmt.Sprintf("%s price in USD: %.2f", crypto, price)
		err = ProduceMessage(topic, message)
		if err != nil {
			fmt.Printf("Error producing message for %s: %s\n", crypto, err)
		} else {
			fmt.Printf("%s price fetched and produced to Kafka: %.2f USD\n", crypto, price)
		}
	}

	return nil
}
