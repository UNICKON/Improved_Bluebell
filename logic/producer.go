package logic

import (
	"awesomeProject/dao/kafka"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaProducer 发送消息到Kafka（统一连接）
func KafkaProducer(topic string, value interface{}) error {
	msgBytes, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("Kafka marshal failed", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	zap.L().Info(msg.Topic, zap.Any("value", value))

	_, _, err = kafka.Producer.SendMessage(msg)

	if err != nil {
		zap.L().Error("Kafka send failed", zap.Error(err))
	} else {
		zap.L().Info("Kafka send success")
	}
	return err
}
