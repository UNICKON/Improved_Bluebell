package logic

import (
	"awesomeProject/dao/kafka"
	"awesomeProject/dao/mysql"
	"database/sql"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// 定期扫描本地消息表，重试发送失败的消息
func RetryUnsentMessages(db *sql.DB, limit int) {
	msgs, err := mysql.GetUnsentMessages(limit)
	if err != nil {
		zap.L().Error("get unsent messages failed", zap.Error(err))
		return
	}
	for _, m := range msgs {
		data, _ := json.Marshal(m)
		_, _, err := kafka.Producer.SendMessage(&sarama.ProducerMessage{
			Topic: "feed_event",
			Value: sarama.ByteEncoder(data),
		})
		if err != nil {
			_ = mysql.MarkMessageFailed(m["MsgID"].(string))
			zap.L().Error("retry send mq failed", zap.Error(err))
			continue
		}
		_ = mysql.MarkMessageSent(m["MsgID"].(string))
	}
}
