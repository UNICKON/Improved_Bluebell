package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/models"
	snowflake "awesomeProject/pkg"
	"awesomeProject/pkg/errorcode"
	"awesomeProject/pkg/jwt"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1.一致性检验
	exist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if exist {
		return errorcode.ErrorUserExist
	}
	//2.生成ID
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.入库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (aToken, rToken string, err error) {
	exit, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return "", "", err
	}
	if !exit {
		return "", "", errorcode.ErrorUserExist
	}

	var pwd string
	pwd, err = mysql.GetUserByUsername(p)
	if err != nil {
		return "", "", err
	}

	p.Password = mysql.EncryptPassword(p.Password)
	if pwd != p.Password {
		return "", "", errorcode.ErrorWrongPassword
	}
	//此时p中获得了userid
	aToken, rToken, err = jwt.GenToken(p.UserId, p.Username)
	if err != nil {
		return "", "", err
	}
	return
}
