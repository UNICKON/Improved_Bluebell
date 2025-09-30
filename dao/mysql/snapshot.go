package mysql

import (
	"awesomeProject/models"
)

// SaveHotRankSnapshotToMySQL 保存热榜快照到 MySQL
func SaveHotRankSnapshotToMySQL(snapshot []models.HotRankSnapshot) error {
	if len(snapshot) == 0 {
		return nil
	}
	query := `
	INSERT INTO hot_rank_snapshot (post_id, score, snapshot_time)
	VALUES (:post_id, :score, :snapshot_time)
	ON DUPLICATE KEY UPDATE score = VALUES(score), snapshot_time = VALUES(snapshot_time);
	`
	_, err := db.NamedExec(query, snapshot)
	return err
}

// GetHotRankSnapshotFromMySQL 查询历史快照
func GetHotRankSnapshotFromMySQL(snapshotTime string, limit int) ([]models.HotRankSnapshot, error) {
	var snapshots []models.HotRankSnapshot
	query := `SELECT post_id, score, snapshot_time FROM hot_rank_snapshot WHERE snapshot_time = ? ORDER BY score DESC LIMIT ?`
	err := db.Select(&snapshots, query, snapshotTime, limit)
	return snapshots, err
}
