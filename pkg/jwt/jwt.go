package jwt

import (
	"awesomeProject/pkg/errorcode"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

const TokenExpireDuration = time.Hour * 24 * 360

//const TokenExpireDuration = time.Second * 30

var mySecret = []byte(viper.GetString("jwt_secret"))

func GenToken(userId int64, username string) (aToken, rToken string, err error) {
	c := MyClaims{
		userId,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 360).Unix(),
		Issuer:    "bluebell",
	}).SignedString(mySecret)
	return
}
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errorcode.ErrorInvalidToken

}

//	func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
//		if _, err = jwt.Parse(rToken, keyFunc); err != nil {
//			return "", "", err
//		}
//		// 从旧access token中解析出claims数据
//		var claims MyClaims
//		_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
//		v, _ := err.(jwt.ValidationError)
//
//		// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
//		if v.Errors == jwt.ValidationErrorExpired {
//			return GenToken(claims.UserId, claims.Username)
//		}
//		return
//	}
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回

	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧access token中解析出claims数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)

	if err == nil {
		return "", "", errors.New("aToken still valid")
	}

	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserId, claims.Username)
	}
	return
}
