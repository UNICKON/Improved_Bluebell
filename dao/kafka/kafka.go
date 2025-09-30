package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	Config   *sarama.Config
	Client   sarama.Client
)

// InitKafka 使用 Viper 初始化 Kafka
func InitKafka() error {
	// 读取配置
	brokers := viper.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		zap.L().Error("No Kafka brokers configured")
		return fmt.Errorf("no Kafka brokers configured")
	}

	Config = sarama.NewConfig()
	Config.Producer.Return.Successes = viper.GetBool("kafka.producer.return_successes")
	// 可根据需要设置更多配置，比如 Config.Consumer.Return.Errors

	var err error

	Client, err = sarama.NewClient(brokers, Config)
	if err != nil {
		zap.L().Error("Kafka client init failed", zap.Error(err))
		return err
	}

	Producer, err = sarama.NewSyncProducer(brokers, Config)
	if err != nil {
		zap.L().Error("Kafka producer init failed", zap.Error(err))
		return err
	}

	Consumer, err = sarama.NewConsumer(brokers, Config)
	if err != nil {
		zap.L().Error("Kafka consumer init failed", zap.Error(err))
		return err
	}

	//brokers = []string{"127.0.0.1:9092"}
	//topic := "vote_event"
	//
	//config := sarama.NewConfig()
	//config.Producer.Return.Successes = true
	//
	//admin, err := sarama.NewClusterAdmin(brokers, config)
	//if err != nil {
	//	panic(fmt.Sprintf("Failed to create admin: %v", err))
	//}
	//defer admin.Close()
	//
	//err = admin.CreateTopic(topic, &sarama.TopicDetail{
	//	NumPartitions:     1,
	//	ReplicationFactor: 1,
	//}, false)
	//if err != nil {
	//	zap.L().Error("Failed to create topic", zap.String("topic", topic), zap.Error(err))
	//} else {
	//	zap.L().Info("Created topic", zap.String("topic", topic))
	//}

	zap.L().Info("Kafka initialized successfully", zap.Strings("brokers", brokers))
	return nil

}

// CloseKafka 关闭 Kafka 连接
func CloseKafka() {
	if Producer != nil {
		_ = Producer.Close()
	}
	if Consumer != nil {
		_ = Consumer.Close()
	}
}
