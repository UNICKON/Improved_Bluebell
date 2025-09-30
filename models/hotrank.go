package models

import "time"

// HotRankSnapshot 用于热榜快照的结构体
// 用于 Redis/MySQL 快照存储
// post_id: 帖子ID, score: 热度分数, snapshot_time: 快照时间
// 可扩展字段如 community_id、title 等

type HotRankSnapshot struct {
	PostID      int64     `json:"post_id" db:"post_id"`
	Score       float64   `json:"score" db:"score"`
	SnapshotTime time.Time `json:"snapshot_time" db:"snapshot_time"`
}
