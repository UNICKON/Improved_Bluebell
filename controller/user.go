package controller

import (
	"awesomeProject/logic"
	"awesomeProject/models"
	"awesomeProject/pkg/errorcode"
	"awesomeProject/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("ParamSignUp ShouldBindJSON", zap.Error(err))
		ResponseErrorwithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	//2.业务处理
	err := logic.SignUp(&p)
	if err == errorcode.ErrorUserExist {
		ResponseErrorwithMsg(c, CodeUserAlreadyExists, err.Error())
		return
	}
	if err != nil {
		zap.L().Error("SignUp error", zap.Error(err))
		ResponseErrorwithMsg(c, CodeSignUpFail, err.Error())
		return
	}
	//3.返回响应
	ResponseSuccess(c, "Sign up success")
}

func LoginHandler(c *gin.Context) {
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("ParamSignUp ShouldBindJSON", zap.Error(err))
		ResponseSuccess(c, CodeInvalidParams)
		return
	}
	aToken, rToken, err := logic.Login(&p)
	if err == errorcode.ErrorUserNotExist {
		ResponseErrorwithMsg(c, CodeUserNotExists, err.Error())
		return
	}
	if err == errorcode.ErrorUserExist {
		ResponseErrorwithMsg(c, CodeInvalidPassword, err.Error())
		return
	}
	if err != nil {
		zap.L().Error("Login error", zap.Error(err))
		ResponseErrorwithMsg(c, CodeLoginFail, err.Error())
		return
	}
	ResponseSuccess(c, gin.H{
		"refreshToken": rToken,
		"accessToken":  aToken,
		"userId":       p.UserId,
		"userName":     p.Username,
	})
}

func RefreshTokenHandler(c *gin.Context) {

	rToken := c.Query("refreshToken")
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		zap.L().Error("No Authorization Header")
		ResponseError(c, CodeNeedAuth)
		c.Abort()
		return
	}

	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		zap.L().Error("Invalid Authorization Header")
		ResponseError(c, CodeInvalidToken)
		c.Abort()
		return
	}

	newAToken, newRToken, err := jwt.RefreshToken(parts[1], rToken)
	if err != nil {
		zap.L().Error("RefreshToken error", zap.Error(err))
		ResponseErrorwithMsg(c, CodeInvalidToken, err.Error())
		return
	}
	ResponseSuccess(c, gin.H{
		"refreshToken": newAToken,
		"accessToken":  newRToken,
	})
}
