package logic

import (
	"awesomeProject/dao/kafka"
	"awesomeProject/dao/redis"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"strconv"
)

const FanoutThreshold = 0

// 处理的消息结构
func FeedConsumer() {
	partitionConsumer, err := kafka.Consumer.ConsumePartition("feed_event", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("Kafka partition consumer failed", zap.Error(err))
		return
	}
	defer partitionConsumer.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient("feed_group", kafka.Client)
	if err != nil {
		zap.L().Error("Kafka offset manager failed", zap.Error(err))
		return
	}
	defer offsetManager.Close()

	partitionOffsetManager, err := offsetManager.ManagePartition("feed_event", 0)
	if err != nil {
		zap.L().Error("Kafka partition offset manager failed", zap.Error(err))
		return
	}
	defer partitionOffsetManager.Close()

	for msg := range partitionConsumer.Messages() {
		var m Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			zap.L().Error("feed unmarshal failed", zap.Error(err))
			// ❌ 不提交 offset，消息会被重新消费
			continue
		}

		fansCount, _ := redis.GetFansCount(m.UserID)
		processErr := false
		if fansCount < FanoutThreshold {
			fans, _ := redis.GetCachedFans(m.UserID)
			// 写扩散
			zap.L().Info("write to fans feed:", zap.Any("fans", fans))
			for _, fanID := range fans {
				if err := redis.PushFeed(fanID, m); err != nil {
					processErr = true
					break
				}
			}
		} else {
			// 读扩散，大V消息写入 bigv:feed:<id>
			zap.L().Info("write to bigv feed:", zap.Any("bigv", m.UserID))
			bigvKey := "bigv:feed:" + strconv.FormatInt(m.UserID, 10)
			if err := redis.PushFeedRaw(bigvKey, m); err != nil {
				processErr = true
			}
		}

		if !processErr {
			zap.L().Info("write feed success")
			// ✅ 手动 ack（提交 offset）
			partitionOffsetManager.MarkOffset(msg.Offset+1, "")
			zap.L().Info("feed processed success",
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset))
		} else {
			// ❌ 失败，不提交 offset，消息会被再次消费
			zap.L().Error("feed process failed, not ack", zap.Any("msg", m))
		}
	}
}
