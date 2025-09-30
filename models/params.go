package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserId   int64  `json:"user_id"`
}

type ParamPost struct {
	Title       string `json:"title" db:"title" binding:"required"`
	Content     string `json:"content" db:"content" binding:"required"`
	CommunityID int64  `json:"community_id" db:"community_id" binding:"required"`
	PostID      int64  `json:"post_id" db:"post_id"`
	AuthorId    int64  `json:"author_id" db:"author_id"`
}

type ParamVote struct {
	PostID    int64 `json:"post_id" binding:"required"`
	Direction int8  `json:"direction" binding:"required,oneof=0 1 -1"`
}

// ParamPostList 获取帖子列表query 参数
type ParamPostList struct {
	Search      string `json:"search" form:"search"`               // 关键字搜索
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}
