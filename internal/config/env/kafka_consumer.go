package env

import (
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"strings"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokerStr := os.Getenv(brokersEnvName)
	if len(brokerStr) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", brokersEnvName)
	}

	brokers := strings.Split(brokerStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", groupIDEnvName)
	}

	return &kafkaConsumerConfig{brokers: brokers, groupID: groupID}, nil
}

func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
