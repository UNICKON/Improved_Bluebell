package mysql

import (
	"awesomeProject/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

func CreatePost(p *models.ParamPost) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr,
		p.PostID, p.Title, p.Content, p.AuthorId, p.CommunityID)
	if err != nil {
		zap.L().Error("CreatePost failed", zap.Error(err))
		return
	}
	return
}

func GetPostByID(idStr string) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time, update_time
	from post
	where post_id = ?`
	err = db.Get(post, sqlStr, idStr)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	zap.L().Info("query post success", zap.Any("post", post))
	return
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

func GetPostTotalCount() (count int64, err error) {
	sqlStr := `select count(post_id) from post`
	err = db.Get(&count, sqlStr)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

// GetCommunityPostTotalCount 根据社区id查询数据库帖子总数
func GetCommunityPostTotalCount(communityID int64) (count int64, err error) {
	sqlStr := `select count(post_id) from post where community_id = ?`
	err = db.Get(&count, sqlStr, communityID)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}
