package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id" db:"post_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"-" db:"create_time"`
	UpdateTime  time.Time `json:"-" db:"update_time"`
}

type PostDetail struct {
	*Post
	*CommunityDetail `json:"community"` // 嵌入社区信息
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"` // 投票数量
}

type Page struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type PostDetailRes struct {
	Page Page          `json:"page"`
	List []*PostDetail `json:"list"`
}
