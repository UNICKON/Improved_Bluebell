package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
)

// 关注
func Follow(userID, followID int64) error {
	// 1. MySQL写入
	err := mysql.AddFollow(userID, followID)
	if err != nil {
		return err
	}
	// 2. Redis缓存更新
	follows, _ := mysql.GetFollows(userID)
	_ = redis.CacheFollows(userID, follows)
	fans, _ := mysql.GetFans(followID)
	_ = redis.CacheFans(followID, fans)
	return nil
}

// 取消关注
func Unfollow(userID, followID int64) error {
	// 1. MySQL删除
	err := mysql.RemoveFollow(userID, followID)
	if err != nil {
		return err
	}
	// 2. Redis缓存更新
	follows, _ := mysql.GetFollows(userID)
	_ = redis.CacheFollows(userID, follows)
	fans, _ := mysql.GetFans(followID)
	_ = redis.CacheFans(followID, fans)
	return nil
}
