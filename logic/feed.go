package logic

import (
	"awesomeProject/dao/kafka"
	"awesomeProject/dao/mysql"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// 消息结构体
type Message struct {
	MsgID   string
	UserID  int64
	Content string
	Time    int64
}

// 用户发消息，先写本地消息表，再推送到MQ
func PublishMessage(msg Message) error {
	// 1. 写本地消息表
	err := mysql.SaveLocalMessage(msg.MsgID, msg.UserID, msg.Content, msg.Time)
	if err != nil {
		zap.L().Error("save local message failed", zap.Error(err))
		return err
	}
	// 2. 发送到MQ
	data, _ := json.Marshal(msg)
	_, _, err = kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "feed_event",
		Value: sarama.ByteEncoder(data),
	})
	zap.L().Info("send message to kafka", zap.Any("msg", msg))
	if err != nil {
		_ = mysql.MarkMessageFailed(msg.MsgID)
		zap.L().Error("send mq failed", zap.Error(err))
		return err
	}
	_ = mysql.MarkMessageSent(msg.MsgID)
	return nil
}
