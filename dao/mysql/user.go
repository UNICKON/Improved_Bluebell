package mysql

import (
	"awesomeProject/models"
	"crypto/md5"
	"encoding/hex"
)

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(*) from user where username= ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return false, err
	}

	return count > 0, nil
}

func QueryUserByID() {

}

// InsertUser 插入用户记录
func InsertUser(user *models.User) (err error) {
	//执行sql
	user.Password = EncryptPassword(user.Password)
	sqlStr := `insert into user (user_id, username,password) values (?, ?, ?)`

	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)

	return
}

func GetUserByUsername(user *models.ParamLogin) (pwd string, err error) {
	sqlStr := `SELECT user_id, password FROM user WHERE username = ?`

	// 用于接收查询结果
	var userID int64
	err = db.QueryRow(sqlStr, user.Username).Scan(&userID, &pwd)

	if err != nil {
		return "", err
	}
	// 赋值 user_id
	user.UserId = userID
	return
}

func GetUserByID(idStr string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, idStr)
	return
}

func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
