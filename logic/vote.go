package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
	"awesomeProject/models"
	"awesomeProject/pkg/errorcode"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// 投票操作，保证MySQL和Redis一致性，热榜快照记录
func Vote(userID int64, postID string, direction float64) (err error) {
	lockKey := "lock:vote:" + postID
	locked, err := redis.Lock(lockKey, 3*time.Second)
	if err != nil || !locked {
		zap.L().Warn("vote lock failed", zap.String("key", lockKey), zap.Error(err))
		return errorcode.ErrorServerBusy
	}
	defer redis.Unlock(lockKey)

	// 1. Redis投票，优先写缓存
	err = redis.VoteForPost(strconv.FormatInt(userID, 10), postID, direction)
	if err == errorcode.ErrorVoteTimeExpired {
		// 限流检查
		legal, err := redis.Checklegal(userID)
		if !legal {
			zap.L().Info("illegal vote from user", zap.Int64("user_id", userID))
			return nil
		}
		// MySQL读入点赞数据
		intPostID, _ := strconv.ParseInt(postID, 10, 64)
		votes, err := mysql.GetVotes(intPostID)
		if err != nil {
			zap.L().Error("GetVotes", zap.Error(err))
			return err
		}
		// 写入Redis缓存
		err = redis.StoreVotes(votes, postID)
		// 再次投票
		err = redis.VoteForPost(strconv.FormatInt(userID, 10), postID, direction)
		if errors.Is(err, errorcode.ErrorAlreadyVote) {
			zap.L().Info("vote already vote", zap.Int64("user_id", userID))
			return nil
		}
		if err != nil {
			zap.L().Error("VoteForPost", zap.Error(err))
			return err
		}
	} else if err != nil {
		zap.L().Error("VoteForPost", zap.Error(err))
		return err
	}

	// 2. 异步写MySQL，发送Kafka消息
	voteRecord := models.Vote{
		PostID: func() int64 { id, _ := strconv.ParseInt(postID, 10, 64); return id }(),
		UserID: userID,
		Vote:   int(direction),
	}
	go func() {
		_ = KafkaProducer("vote_event", voteRecord)
	}()

	// 3. 热榜快照记录（每10分钟一次，可用定时任务触发）
	if time.Now().Minute()%10 == 0 {
		go SaveHotRankSnapshot()
	}

	return nil
}

// SaveHotRankSnapshot 保存当前热榜快照到Redis和MySQL
func SaveHotRankSnapshot() {
	lockKey := "lock:hotrank:snapshot"
	locked, err := redis.Lock(lockKey, 10*time.Second)
	if err != nil || !locked {
		zap.L().Warn("hotrank snapshot lock failed", zap.String("key", lockKey), zap.Error(err))
		return
	}
	defer redis.Unlock(lockKey)

	// 1. 获取当前热榜（如ZSet）
	hotRank, err := redis.GetHotRankSnapshot("hot_rank_zset", 100)
	if err != nil {
		zap.L().Error("GetHotRankSnapshot", zap.Error(err))
		return
	}
	// 2. 保存到Redis快照
	snapshotKey := "hot_rank_snapshot:" + time.Now().Format("200601021504")
	postMap := make(map[string]float64)
	var mysqlSnapshot []models.HotRankSnapshot
	for _, z := range hotRank {
		postMap[z.Member.(string)] = z.Score
		postID, _ := strconv.ParseInt(z.Member.(string), 10, 64)
		mysqlSnapshot = append(mysqlSnapshot, models.HotRankSnapshot{
			PostID:       postID,
			Score:        z.Score,
			SnapshotTime: time.Now(),
		})
	}
	_ = redis.SaveHotRankSnapshot(snapshotKey, postMap)
	// 3. 保存到MySQL快照表
	_ = mysql.SaveHotRankSnapshotToMySQL(mysqlSnapshot)
}
