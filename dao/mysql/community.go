package mysql

import (
	"awesomeProject/models"
	"awesomeProject/pkg/errorcode"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

var (
	ErrorInvalidID   = errors.New("Invalid ID")
	ErrorQueryFailed = errors.New("Query failed")
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "Select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}

func GetCommunityDetailByID(idStr string) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
	from community
	where community_id = ?`
	err = db.Get(community, sqlStr, idStr)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func GetCommunityNameByID(idStr string) (community *models.Community, err error) {
	community = new(models.Community)
	sqlStr := `select community_id, community_name
	from community
	where community_id = ?`
	err = db.Get(community, sqlStr, idStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorcode.ErrorInvalidID

		}
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		return nil, errorcode.ErrorInvalidID
	}
	return
}

// GetCommunityByID 根据ID查询分类社区详情
func GetCommunityByID(id int64) (*models.CommunityDetail, error) {
	community := new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
	from community
	where community_id = ?`
	err := db.Get(community, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows { // 查询为空
			return nil, errorcode.ErrorInvalidID
		}
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		return nil, errorcode.ErrorInvalidID
	}
	return &models.CommunityDetail{
		ID:         community.ID,
		Name:       community.Name,
		Intro:      community.Intro,
		CreateTime: community.CreateTime,
	}, err
}
