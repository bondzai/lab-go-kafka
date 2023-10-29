package constants

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	KafkaCryptoPriceTopic = os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")
)

func GetKafkaConfig() kafka.ConfigMap {
	return kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("CLOUDKARAFKA_BROKERS"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     os.Getenv("CLOUDKARAFKA_USERNAME"),
		"sasl.password":     os.Getenv("CLOUDKARAFKA_PASSWORD"),
	}
}
