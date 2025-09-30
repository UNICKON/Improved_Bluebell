package redis

import (
	"github.com/go-redis/redis"
	"time"
)

// SaveHotRankSnapshot 保存当前热榜快照到 Redis
func SaveHotRankSnapshot(snapshotKey string, hotRank map[string]float64) error {
	pipe := rdb.TxPipeline()
	for postID, score := range hotRank {
		pipe.ZAdd(snapshotKey, redis.Z{Score: score, Member: postID})
	}
	pipe.Expire(snapshotKey, 24*time.Hour) // 快照保留24小时
	_, err := pipe.Exec()
	return err
}

// GetHotRankSnapshot 获取指定快照的热榜数据
func GetHotRankSnapshot(snapshotKey string, count int64) ([]redis.Z, error) {
	return rdb.ZRevRangeWithScores(snapshotKey, 0, count-1).Result()
}
