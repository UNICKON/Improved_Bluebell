package mysql

// 关注关系表CRUD
func AddFollow(userID, followID int64) error {
	_, err := db.Exec("INSERT INTO user_follow (user_id, follow_id) VALUES (?, ?)", userID, followID)
	return err
}

func RemoveFollow(userID, followID int64) error {
	_, err := db.Exec("DELETE FROM user_follow WHERE user_id=? AND follow_id=?", userID, followID)
	return err
}

func GetFollows(userID int64) ([]int64, error) {
	rows, err := db.Query("SELECT follow_id FROM user_follow WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	var follows []int64
	for rows.Next() {
		var fid int64
		_ = rows.Scan(&fid)
		follows = append(follows, fid)
	}
	return follows, nil
}

func GetFans(followID int64) ([]int64, error) {
	rows, err := db.Query("SELECT user_id FROM user_follow WHERE follow_id=?", followID)
	if err != nil {
		return nil, err
	}
	var fans []int64
	for rows.Next() {
		var uid int64
		_ = rows.Scan(&uid)
		fans = append(fans, uid)
	}
	return fans, nil
}
