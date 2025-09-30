package controller

import (
	"awesomeProject/logic"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Vote struct {
	PostID    string `json:"post_id,string"`
	Direction string `json:"direction"`
}

func VoteHandler(c *gin.Context) {
	vote := new(models.ParamVote)
	if err := c.ShouldBindJSON(&vote); err != nil {
		zap.L().Error("c.ShouldBindJSON(vote) failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}

	if err = logic.Vote(userID, strconv.Itoa(int(vote.PostID)), float64(vote.Direction)); err != nil {
		zap.L().Error("logic.Vote(userID,vote.Direction) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return

	}

	ResponseSuccess(c, nil)

}
