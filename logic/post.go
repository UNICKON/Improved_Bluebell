package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
	"awesomeProject/models"
	snowflake "awesomeProject/pkg"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func CreatePost(p *models.ParamPost) (err error) {
	postID := snowflake.GenID()
	p.PostID = postID

	if err = mysql.CreatePost(p); err != nil {
		zap.L().Error("CreatePost failed", zap.Error(err))
		return err
	}

	community, err := mysql.GetCommunityNameByID(fmt.Sprint(p.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}
	// redis存储帖子信息
	if err := redis.CreatePost(
		p.PostID,
		p.AuthorId,
		p.Title,
		TruncateByWords(p.Content, 120),
		community.ID); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return
}
func GetPost(postID string) (res *models.PostDetail, err error) {
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed", zap.String("post_id", postID), zap.Error(err))
		return nil, err
	}
	user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))

	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
		return
	}
	community, err := mysql.GetCommunityDetailByID(fmt.Sprint(post.CommunityID))

	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
		return
	}
	// 根据帖子id查询帖子的投票数
	voteNum, err := redis.GetPostVoteNum(postID)

	// 接口数据拼接
	res = &models.PostDetail{
		Post:            post,
		CommunityDetail: community,
		AuthorName:      user.Username,
		VoteNum:         voteNum,
	}
	fmt.Printf("%+v\n", res)

	return
}

// GetPostListNew 将两个查询帖子列表逻辑合二为一的函数
func GetPostListNew(p *models.ParamPostList) (data *models.PostDetailRes, err error) {
	// 根据请求参数的不同,执行不同的业务逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func GetPostList2(p *models.ParamPostList) (res *models.PostDetailRes, err error) {
	res = new(models.PostDetailRes)
	// 从mysql获取帖子列表总数
	total, err := mysql.GetPostTotalCount()
	if err != nil {
		zap.L().Error("mysql.GetPostTotalCount failed", zap.Error(err))
		return nil, err
	}
	res.Page.Total = total
	// 1、根据参数中的排序规则去redis查询id列表 #page limit
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInorder failed", zap.Error(err))
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder return empty array")
		return nil, nil
	}
	// 2、提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData failed", zap.Error(err))
		return nil, err
	}

	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs failed", zap.Error(err))
		return nil, err
	}

	//4.组合数据
	res.Page.Page = p.Page
	res.Page.Size = p.Page
	res.List = make([]*models.PostDetail, 0, len(ids))
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		user, err := mysql.GetUserByID(strconv.Itoa(int(post.AuthorId)))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityDetailByID(fmt.Sprint(post.CommunityID))
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Error(err))
			return nil, err
		}
		// 接口数据拼接
		postDetail := &models.PostDetail{
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.Username,
		}
		res.List = append(res.List, postDetail)
	}
	return
}

// GetCommunityPostList 根据社区id去查询帖子列表
func GetCommunityPostList(p *models.ParamPostList) (*models.PostDetailRes, error) {
	var res models.PostDetailRes
	// 从mysql获取该社区下帖子列表总数
	total, err := mysql.GetCommunityPostTotalCount(p.CommunityID)
	if err != nil {
		return nil, err
	}
	res.Page.Total = total
	// 1、根据参数中的排序规则去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) return 0 data")
		return &res, nil
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 2、提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Page = p.Page
	res.Page.Size = p.Size
	res.List = make([]*models.PostDetail, 0, len(posts))
	// 4、根据社区id查询社区详细信息
	// 为了减少数据库的查询次数，这里将社区信息提前查询出来
	community, err := mysql.GetCommunityByID(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Int("id", int(p.CommunityID)),
			zap.Error(err))
		community = nil
	}
	for idx, post := range posts {
		// 过滤掉不属于该社区的帖子
		if post.CommunityID != p.CommunityID {
			continue
		}
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(strconv.FormatInt(post.AuthorId, 10))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", uint64(post.AuthorId)),
				zap.Error(err))
			user = nil
		}
		// 接口数据拼接
		postDetail := &models.PostDetail{
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.Username,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}
