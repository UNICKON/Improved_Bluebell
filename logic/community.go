package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id string) (CommunityDetail *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
