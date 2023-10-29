package constants

import "os"

var (
	KafkaBootstrapServers = os.Getenv("CLOUDKARAFKA_BROKERS")
	KafkaUsername         = os.Getenv("CLOUDKARAFKA_USERNAME")
	KafkaPassword         = os.Getenv("CLOUDKARAFKA_PASSWORD")
	KafkaSecureProtocol   = "SASL_SSL"
	KafkaSaslMech         = "SCRAM-SHA-512"
	KafkaCryptoPriceTopic = os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")
)
