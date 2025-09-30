package controller

import (
	"awesomeProject/logic"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func CreatePostHandler(c *gin.Context) {
	var p models.ParamPost
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("binding", zap.Any("err", err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	id, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("Error getting current user id", zap.Error(err))
		ResponseError(c, CodeNotLogin)
	}
	p.AuthorId = id

	err = logic.CreatePost(&p)
	if err != nil {
		zap.L().Error("Error creating post", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 发帖成功后自动推送feed
	msg := logic.Message{
		MsgID:   strconv.FormatInt(p.PostID, 10),
		UserID:  id,
		Content: p.Title + "\n" + p.Content,
		Time:    time.Now().Unix(),
	}
	_ = logic.PublishMessage(msg)
	ResponseSuccess(c, nil)
}

// PostDetailHandler 帖子详情
func PostDetailHandler(c *gin.Context) {
	postId := c.Param("id")
	post, err := logic.GetPost(postId)
	if err != nil {
		zap.L().Error("logic.GetPost(postID) failed", zap.String("postId", postId), zap.Error(err))
	}

	ResponseSuccess(c, post)
}

func PostList2Handler(c *gin.Context) {
	// GET请求参数(query string)： /api/v1/posts2?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostList2Handler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取数据
	data, err := logic.GetPostListNew(p) // 更新：合二为一
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
