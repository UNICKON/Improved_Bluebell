package mysql

import (
	"awesomeProject/models"
)

func StoreExpirePostToMySQL(votes []models.Vote) error {
	if len(votes) == 0 {
		return nil
	}

	// 构造 SQL 语句
	query := `
	INSERT INTO vote (post_id, user_id, vote)
	VALUES (:post_id, :user_id, :vote)
	ON DUPLICATE KEY UPDATE vote = VALUES(vote),
	UPDATED_AT = NOW();
	`

	// 使用 NamedExec 批量执行
	_, err := db.NamedExec(query, votes)
	if err != nil {
		return err
	}
	return nil
}

func GetVotes(postID int64) ([]models.Vote, error) {
	var vote []models.Vote
	query := `SELECT post_id, user_id, vote FROM vote WHERE post_id = ?`
	err := db.Select(&vote, query, postID)
	if err != nil {
		return nil, err
	}
	return vote, nil
}
