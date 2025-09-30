package logic

import (
	"awesomeProject/dao/redis"
	"encoding/json"
	"sort"
	"strconv"
	"time"
)

// 普通用户拉取feed
func GetUserFeed(userID int64) ([]map[string]interface{}, error) {
	feedStrs, err := redis.GetFeed(userID, 50)
	if err != nil {
		return nil, err
	}
	var feeds []map[string]interface{}
	for _, s := range feedStrs {
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(s), &m)
		feeds = append(feeds, m)
	}
	return feeds, nil
}

// 大V粉丝拉取feed（聚合）
func GetBigVFeed(userID int64) ([]map[string]interface{}, error) {
	follows, _ := redis.GetCachedFollows(userID)
	var allFeeds []map[string]interface{}
	for _, fid := range follows {
		feedStrs, _ := redis.GetFeed(fid, 20)
		for _, s := range feedStrs {
			var m map[string]interface{}
			_ = json.Unmarshal([]byte(s), &m)
			allFeeds = append(allFeeds, m)
		}
	}
	// 可按时间排序
	return allFeeds, nil
}

// 用户拉取关注消息（混合策略，修正大Vfeed处理）
func GetFollowFeed(userID int64) ([]map[string]interface{}, error) {
	var allFeeds []map[string]interface{}
	// 1. 拉取自己的feed流
	selfFeeds, err := redis.GetFeed(userID, 50)
	if err == nil {
		for _, s := range selfFeeds {
			var m map[string]interface{}
			_ = json.Unmarshal([]byte(s), &m)
			allFeeds = append(allFeeds, m)
		}
	}
	// 2. 拉取关注大V的 bigv feed 流
	follows, _ := redis.GetCachedFollows(userID)
	for _, fid := range follows {
		fansCount, _ := redis.GetFansCount(fid)
		if fansCount >= FanoutThreshold {
			bigvKey := "bigv:feed:" + strconv.FormatInt(fid, 10)
			feedStrs, _ := redis.GetRawFeed(bigvKey, 20) // 新增GetRawFeed方法，按key拉取
			for _, s := range feedStrs {
				var m map[string]interface{}
				_ = json.Unmarshal([]byte(s), &m)
				allFeeds = append(allFeeds, m)
			}
		}
	}
	// 3. 只保留24小时内发布的消息
	cutoff := time.Now().Add(-24 * time.Hour).Unix()
	var recentFeeds []map[string]interface{}
	for _, m := range allFeeds {
		if ts, ok := m["Time"].(float64); ok && int64(ts) >= cutoff {
			recentFeeds = append(recentFeeds, m)
		}
	}
	// 4. 按时间降序排序
	sort.Slice(recentFeeds, func(i, j int) bool {
		return recentFeeds[i]["Time"].(float64) > recentFeeds[j]["Time"].(float64)
	})
	return recentFeeds, nil
}
