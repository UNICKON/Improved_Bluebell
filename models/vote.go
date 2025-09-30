package models

type Vote struct {
	PostID int64 `db:"post_id"`
	UserID int64 `db:"user_id"`
	Vote   int   `db:"vote"` // 1 表示赞成，-1 反对，可扩展
}
