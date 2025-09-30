package logic

import (
	"awesomeProject/dao/kafka"
	"awesomeProject/dao/mysql"
	"awesomeProject/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaConsumer 消费投票消息并写入MySQL（统一连接）
// KafkaConsumer 消费投票消息并写入MySQL（统一连接）
func KafkaConsumer(topic string) {
	// 1. 创建 offsetManager
	offsetManager, err := sarama.NewOffsetManagerFromClient("vote_group", kafka.Client)
	if err != nil {
		zap.L().Error("Kafka offset manager failed", zap.Error(err))
		return
	}
	defer offsetManager.Close()

	// 2. 针对 topic 的 0 号分区创建 partitionOffsetManager
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, 0)
	if err != nil {
		zap.L().Error("Kafka partition offset manager failed", zap.Error(err))
		return
	}
	defer partitionOffsetManager.Close()

	// 3. 从上次提交的 offset 开始消费（如果没有，就用 Newest）
	initialOffset, _ := partitionOffsetManager.NextOffset()
	if initialOffset == sarama.OffsetNewest || initialOffset < 0 {
		initialOffset = sarama.OffsetNewest
	}

	partitionConsumer, err := kafka.Consumer.ConsumePartition(topic, 0, initialOffset)
	if err != nil {
		zap.L().Error("Kafka partition consumer failed", zap.Error(err))
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var vote models.Vote
		if err := json.Unmarshal(msg.Value, &vote); err != nil {
			zap.L().Error("vote unmarshal failed", zap.Error(err))
			continue
		}

		// 异步写入MySQL
		if err := mysql.StoreExpirePostToMySQL([]models.Vote{vote}); err != nil {
			zap.L().Error("StoreExpirePostToMySQL", zap.Error(err))
			continue // 写库失败，不提交 offset，下次还会消费
		}

		// ✅ 成功处理才提交 offset
		zap.L().Info("consume success", zap.String("topic", topic), zap.Int64("offset", msg.Offset), zap.String("value", string(msg.Value)))
		partitionOffsetManager.MarkOffset(msg.Offset+1, "")
	}
}
