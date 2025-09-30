package mysql

// 本地消息表结构：id, msg_id, user_id, content, time, status
// status: 0-待发送, 1-已发送, 2-失败

func SaveLocalMessage(msgID string, userID int64, content string, ts int64) error {
	_, err := db.Exec("INSERT INTO local_message (msg_id, user_id, content, time, status) VALUES (?, ?, ?, ?, 0)", msgID, userID, content, ts)
	return err
}

func MarkMessageSent(msgID string) error {
	_, err := db.Exec("UPDATE local_message SET status=1 WHERE msg_id=?", msgID)
	return err
}

func MarkMessageFailed(msgID string) error {
	_, err := db.Exec("UPDATE local_message SET status=2 WHERE msg_id=?", msgID)
	return err
}

func GetUnsentMessages(limit int) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT msg_id, user_id, content, time FROM local_message WHERE status=0 LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	var msgs []map[string]interface{}
	for rows.Next() {
		var msgID string
		var userID int64
		var content string
		var ts int64
		_ = rows.Scan(&msgID, &userID, &content, &ts)
		msgs = append(msgs, map[string]interface{}{
			"MsgID":   msgID,
			"UserID":  userID,
			"Content": content,
			"Time":    ts,
		})
	}
	return msgs, nil
}
