package controller

import (
	"awesomeProject/logic"
	"github.com/gin-gonic/gin"
	"time"
)

type FeedPublishRequest struct {
	Content string `json:"content" binding:"required"`
}

// 用户发消息接口（建议在发帖时自动调用）
func PublishMessageHandler(c *gin.Context) {
	var req FeedPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil || userID == 0 {
		ResponseError(c, CodeNotLogin)
		return
	}
	msg := logic.Message{
		UserID:  userID,
		Content: req.Content,
		Time:    time.Now().Unix(),
	}
	if err := logic.PublishMessage(msg); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// 拉取feed接口
func GetFeedHandler(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil || userID == 0 {
		ResponseError(c, CodeNotLogin)
		return
	}
	feeds, _ := logic.GetFollowFeed(userID)
	ResponseSuccess(c, feeds)
}
