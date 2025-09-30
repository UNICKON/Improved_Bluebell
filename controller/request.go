package controller

import (
	"awesomeProject/pkg/errorcode"
	"github.com/gin-gonic/gin"
)

const CtxUserIdKey = "userID"

func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	_userID, ok := c.Get(CtxUserIdKey)
	if !ok {
		err = errorcode.ErrorUserNotLogin
		return
	}
	userID, ok = _userID.(int64)
	if !ok {
		err = errorcode.ErrorUserNotLogin
		return
	}
	return
}
