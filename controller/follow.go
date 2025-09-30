package controller

import (
	"awesomeProject/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FollowRequest struct {
	FollowID int64 `json:"follow_id" binding:"required"`
}

// 关注接口
func FollowHandler(c *gin.Context) {
	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil || userID == 0 {
		ResponseError(c, CodeNotLogin)
		zap.L().Error("FollowHandler", zap.Error(err))
		return
	}
	if err := logic.Follow(userID, req.FollowID); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// 取消关注接口
func UnfollowHandler(c *gin.Context) {
	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil || userID == 0 {
		ResponseError(c, CodeNotLogin)
		return
	}
	if err := logic.Unfollow(userID, req.FollowID); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
