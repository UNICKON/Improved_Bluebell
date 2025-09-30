package redis

import (
	"awesomeProject/models"
	"awesomeProject/pkg/errorcode"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"
)

const (
	ShortWindowLimit = 10000                 // 1分钟最多访问3次
	LongWindowLimit  = 10000                 // 10分钟最多访问10次
	ShortExpire      = 60 * time.Second      // 1分钟
	LongExpire       = 10 * 60 * time.Second // 10分钟

	OneWeekInSeconds          = 7 * 24 * 3600        // 一周的秒数
	OneMonthInSeconds         = 4 * OneWeekInSeconds // 一个月的秒数
	VoteScore         float64 = 432                  // 每一票的值432分
	PostPerAge                = 20                   // 每页显示20条帖子
)

// VoteForPost	为帖子投票
func VoteForPost(userID string, postID string, v float64) (err error) {
	// 1.判断投票限制
	_, err = rdb.ZScore(KeyPostExpire, postID).Result()
	if err == redis.Nil {
		// key 或 member 不存在，说明帖子没有过期记录，禁止投票
		return errorcode.ErrorVoteTimeExpired
	} else if err != nil {
		// 其他 Redis 错误
		return err
	}
	// 2、更新帖子的分数
	// 2和3 需要放到一个pipeline事务中操作
	// 判断是否已经投过票 查当前用户给当前帖子的投票记录
	key := KeyPostVotedZSetPrefix + postID
	ov := rdb.ZScore(key, userID).Val()
	zap.L().Info("VoteForPost", zap.String("key", key), zap.Float64("ov", ov), zap.String("user", userID))
	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if v == ov {
		return errorcode.ErrorAlreadyVote
	}
	var op float64
	if v > ov {
		op = 1
	} else {
		op = -1
	}
	diffAbs := math.Abs(ov - v)                // 计算两次投票的差值
	pipeline := rdb.TxPipeline()               // 事务操作
	incrementScore := VoteScore * diffAbs * op // 计算分数（新增）
	// ZIncrBy 用于将有序集合中的成员分数增加指定数量
	_, err = pipeline.ZIncrBy(KeyPostScoreZSet, incrementScore, postID).Result() // 更新分数
	if err != nil {
		return err
	}
	// 3、记录用户为该帖子投票的数据
	if v == 0 {
		_, err = rdb.ZRem(key, postID).Result()
	} else {
		pipeline.ZAdd(key, redis.Z{ // 记录已投票
			Score:  v, // 赞成票还是反对票
			Member: userID,
		})
	}
	// 4、更新帖子的投票数
	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", int64(op))

	//switch math.Abs(ov) - math.Abs(v) {
	//case 1:
	//	// 取消投票 ov=1/-1 v=0
	//	// 投票数-1
	//	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -1)
	//case 0:
	//	// 反转投票 ov=-1/1 v=1/-1
	//	// 投票数不用更新
	//case -1:
	//	// 新增投票 ov=0 v=1/-1
	//	// 投票数+1
	//	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 1)
	//default:
	//	// 已经投过票了
	//	return ErrorVoted
	//}
	_, err = pipeline.Exec()
	return err
}

// 返回vote id对
func GetVoteData(postIDs []int64) ([]models.Vote, error) {
	var result []models.Vote

	for _, postID := range postIDs {
		key := KeyPostVotedZSetPrefix + strconv.Itoa(int(postID))
		// 获取全部投票数据
		zMembers, err := rdb.ZRangeWithScores(key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		for _, z := range zMembers {
			UserID, _ := strconv.Atoi(z.Member.(string))
			result = append(result, models.Vote{
				PostID: postID,
				UserID: int64(UserID),
				Vote:   int(z.Score),
			})
		}
	}

	return result, nil
}

func DeleteExpiredPostIDs(postIDs []int64) error {
	if len(postIDs) == 0 {
		return nil
	}

	backupKey := KeyPostExpireBackup

	for _, postID := range postIDs {
		zap.L().Info("deleteExpiredPostIDs", zap.Int64("postID", postID))

		pipe := rdb.Pipeline()

		// 构造点赞 key
		member := strconv.FormatInt(postID, 10)
		// 从过期备份集合中删除 postID
		pipe.ZRem(backupKey, member)
		// 删除对应的点赞数据 key
		voteKey := KeyPostVotedZSetPrefix + strconv.FormatInt(postID, 10)
		pipe.Del(voteKey)
		_, err := pipe.Exec()
		if err != nil {
			zap.L().Error("fail Del", zap.Int64("postID", postID), zap.Error(err))
		}
	}

	return nil
}

func StoreVotes(votes []models.Vote, postID string) error {
	pipe := rdb.Pipeline()
	intID, _ := strconv.ParseInt(postID, 10, 64)
	pipe.ZAdd(KeyPostExpire, redis.Z{
		Score:  float64(time.Now().Unix()) + PostExpireDurationSeconds,
		Member: intID,
	})

	for _, v := range votes {

		key := KeyPostVotedZSetPrefix + strconv.FormatInt(v.PostID, 10)

		pipe.ZAdd(key, redis.Z{
			Score:  float64(v.Vote),
			Member: v.UserID,
		})

		zap.L().Info("storeVotes", zap.Int64("postID", v.PostID), zap.Int64("userID", v.UserID))
	}

	_, err := pipe.Exec()
	return err
}

func Checklegal(userID int64) (bool, error) {
	strUserID := strconv.FormatInt(userID, 10)
	shortKey := "limit:1min:" + strUserID
	longKey := "limit:10min:" + strUserID

	pipe := rdb.TxPipeline()

	// 每次访问就 INCR
	shortIncr := pipe.Incr(shortKey)
	pipe.Expire(shortKey, ShortExpire)

	longIncr := pipe.Incr(longKey)
	pipe.Expire(longKey, LongExpire)

	_, err := pipe.Exec()
	if err != nil {
		return false, err
	}

	// 判断是否超过阈值
	if shortIncr.Val() > int64(ShortWindowLimit) || longIncr.Val() > int64(LongWindowLimit) {
		return false, nil // 是非法请求（超出频率限制）
	}

	return true, nil // 合法
}
