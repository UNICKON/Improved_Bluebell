package redis

import (
	"encoding/json"
	"strconv"
)

// 推送消息到用户feed队列（List）
func PushFeed(userID int64, msg interface{}) error {
	key := "user:feed:" + strconv.FormatInt(userID, 10)
	data, _ := json.Marshal(msg)
	return rdb.LPush(key, data).Err()
}

// 获取用户feed队列
func GetFeed(userID int64, count int64) ([]string, error) {
	key := "user:feed:" + strconv.FormatInt(userID, 10)
	return rdb.LRange(key, 0, count-1).Result()
}

// 按指定key获取feed队列
func GetRawFeed(key string, count int64) ([]string, error) {
	return rdb.LRange(key, 0, count-1).Result()
}

// 按指定key推送消息到feed队列
func PushFeedRaw(key string, msg interface{}) error {
	data, _ := json.Marshal(msg)
	return rdb.LPush(key, data).Err()
}
