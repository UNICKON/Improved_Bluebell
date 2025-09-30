package redis

import (
	"fmt"
	"strconv"
)

// Redis缓存关注关系
func CacheFollows(userID int64, follows []int64) error {
	key := "user:follows:" + strconv.FormatInt(userID, 10)
	members := make([]string, len(follows))
	for i, fid := range follows {
		members[i] = strconv.FormatInt(fid, 10)
	}
	return rdb.SAdd(key, members).Err()
}

func GetCachedFollows(userID int64) ([]int64, error) {
	key := "user:follows:" + strconv.FormatInt(userID, 10)
	members, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	var follows []int64
	for _, m := range members {
		fid, _ := strconv.ParseInt(m, 10, 64)
		follows = append(follows, fid)
	}
	return follows, nil
}

func CacheFans(followID int64, fans []int64) error {
	key := "user:fans:" + strconv.FormatInt(followID, 10)
	members := make([]string, len(fans))
	for i, uid := range fans {
		members[i] = strconv.FormatInt(uid, 10)
	}
	return rdb.SAdd(key, members).Err()
}

// GetFansCount 返回粉丝数
func GetFansCount(followID int64) (int64, error) {
	key := "user:fans:" + strconv.FormatInt(followID, 10)
	count, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis SCard failed: %w", err)
	}
	return count, nil
}

func GetCachedFans(followID int64) ([]int64, error) {
	key := "user:fans:" + strconv.FormatInt(followID, 10)
	members, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	var fans []int64
	for _, m := range members {
		uid, _ := strconv.ParseInt(m, 10, 64)
		fans = append(fans, uid)
	}
	return fans, nil
}
