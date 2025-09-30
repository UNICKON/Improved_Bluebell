package redis

import (
	"awesomeProject/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const PostExpireDurationSeconds = 10

func CreatePost(postID, userID int64, title, summary string, CommunityID int64) (err error) {
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + strconv.Itoa(int(postID))
	communityKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(CommunityID))
	postInfo := map[string]interface{}{
		"title":    title,
		"summary":  summary,
		"post:id":  postID,
		"user:id":  userID,
		"time":     now,
		"votes":    1,
		"comments": 0,
	}

	// 事务操作
	pipeline := rdb.TxPipeline()
	// 投票 zSet
	pipeline.ZAdd(votedKey, redis.Z{ // 作者默认投赞成票
		Score:  1,
		Member: userID,
	})
	//pipeline.Expire(votedKey, time.Second*OneMonthInSeconds*6) // 过期时间：6个月
	// 文章 hash
	pipeline.HMSet(KeyPostInfoHashPrefix+strconv.Itoa(int(postID)), postInfo)
	// 添加到分数 ZSet
	pipeline.ZAdd(KeyPostScoreZSet, redis.Z{
		Score:  now + VoteScore,
		Member: postID,
	})
	// 添加到时间 ZSet
	pipeline.ZAdd(KeyPostTimeZSet, redis.Z{
		Score:  now,
		Member: postID,
	})

	pipeline.ZAdd(KeyPostExpire, redis.Z{
		Score:  now + PostExpireDurationSeconds,
		Member: postID,
	})

	// 添加到对应版块 把帖子添加到社区 set
	pipeline.SAdd(communityKey, postID)
	_, err = pipeline.Exec()
	return
}

// getIDsFormKey 按照分数从大到小的顺序查询指定数量的元素
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZRevRange 按照分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 升级版投票列表接口：按创建时间排序 或者 按照 分数排序 (查询出的ids已经根据order从大到小排序)
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	key := KeyPostTimeZSet            // 默认是时间
	if p.Order == models.OrderScore { // 按照分数请求
		key = KeyPostScoreZSet
	}
	// 2.确定查询的索引起始点
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := KeyPostVotedZSetPrefix + id
		// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
		v := rdb.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	// 使用 pipeline一次发送多条命令减少RTT
	//pipeline := client.Pipeline()
	//for _, id := range ids {
	//	key := KeyCommunityPostSetPrefix + id
	//	pipeline.ZCount(key, "1", "1") // ZCount会返回分数在min和max范围内的成员数量
	//}
	//cmders, err := pipeline.Exec()
	//if err != nil {
	//	return nil, err
	//}
	//data = make([]int64, 0, len(cmders))
	//for _, cmder := range cmders {
	//	v := cmder.(*redis.IntCmd).Val()
	//	data = append(data, v)
	//}
	return data, nil
}

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	orderkey := KeyPostTimeZSet       // 默认是时间
	if p.Order == models.OrderScore { // 按照分数请求
		orderkey = KeyPostScoreZSet
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID))

	// 利用缓存key减少zinterstore执行的次数 缓存key
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // 将两个zset函数聚合的时候 求最大值
		}, cKey, orderkey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteNum 根据id查询每篇帖子的投赞成票的数据
func GetPostVoteNum(ids string) (data int64, err error) {
	key := KeyPostVotedZSetPrefix + ids
	// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
	data = rdb.ZCount(key, "1", "1").Val()
	return data, nil
}

//func GetExpiredPostIDs(limit int64) ([]int64, error) {
//	key := KeyPostExpire
//	// 构造范围查询条件
//	opt := redis.ZRangeBy{
//		Min:    "0",
//		Max:    strconv.FormatInt(time.Now().Unix(), 10), // 当前时间戳
//		Offset: 0,
//		Count:  limit,
//	}
//	// 获取已过期的 postID（按 score 排序）
//	postIDStrs, err := rdb.ZRangeByScore(key, opt).Result()
//	if err != nil {
//		return nil, err
//	}
//	if len(postIDStrs) == 0 {
//		zap.L().Info("empty postIDs")
//		return nil, nil
//	}
//	var postIDs []int64
//	for _, str := range postIDStrs {
//		id, err := strconv.ParseInt(str, 10, 64)
//		if err != nil {
//			continue
//		}
//		postIDs = append(postIDs, id)
//	}
//	return postIDs, nil
//}

func GetExpiredPostIDs(limit int64) ([]int64, error) {
	key := KeyPostExpire
	backupKey := KeyPostExpireBackup

	// 构造范围查询条件
	opt := redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(time.Now().Unix(), 10), // 当前时间戳
		Offset: 0,
		Count:  limit,
	}
	// 获取已过期的 postID（按 score 排序）
	postIDStrs, err := rdb.ZRangeByScore(key, opt).Result()
	if err != nil {
		return nil, err
	}

	var postIDs []int64
	if len(postIDStrs) == 0 {
		return postIDs, nil
	}

	// 开启 pipeline 批量处理

	now := float64(time.Now().Unix())

	for _, str := range postIDStrs {
		// 删除旧 key 中的 postID
		pipe := rdb.Pipeline()

		pipe.ZRem(key, str)

		// 添加到备份 key，score 为当前时间
		pipe.ZAdd(backupKey, redis.Z{
			Score:  now,
			Member: str,
		})

		id, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			postIDs = append(postIDs, id)
		}
		_, err = pipe.Exec()

		if err != nil {
			return nil, err
		}
	}

	return postIDs, nil
}
