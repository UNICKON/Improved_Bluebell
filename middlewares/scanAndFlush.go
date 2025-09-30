package middlewares

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
	"awesomeProject/logic"
	"go.uber.org/zap"
	"time"
)

func StartLikeFlusher() {
	go logic.KafkaConsumer("vote_event")
	go logic.FeedConsumer()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		//zap.L().Info("▶️ 点赞数据扫描启动：")
		scanAndFlushLikes()
	}
}

func scanAndFlushLikes() {
	//获得最早的postid
	result, err := redis.GetExpiredPostIDs(100)
	if err != nil {
		zap.L().Error("rdb.ZRangeWithScores() failed", zap.Error(err))
	}

	if len(result) != 0 {

		//获得postid的点赞数据
		res, err := redis.GetVoteData(result)

		zap.L().Info("rdb.ZRevRangeWithScores() result", zap.Any("result", res))

		if err != nil {
			zap.L().Error("redis.GetVoteData() failed", zap.Error(err))
		}
		// 将数据存回 MySQL
		err = mysql.StoreExpirePostToMySQL(res)
		if err != nil {
			zap.L().Error("storePostToMySQL", zap.Error(err))
		}

		//删除对应的redis数据
		err = redis.DeleteExpiredPostIDs(result)
	}
}
