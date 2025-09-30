package controller

type Rescode int

const (
	CodeSuccess Rescode = 1000 + iota
	CodeInvalidParams
	CodeUserNotExists
	CodeUserAlreadyExists
	CodeInvalidPassword
	CodeServerBusy
	CodeLoginFail
	CodeSignUpFail
	CodeNotLogin
	CodeNeedAuth
	CodeInvalidToken
	CodeTokenExpired
)

var codeMsgMap = map[Rescode]string{
	CodeSuccess:           "success",
	CodeInvalidParams:     "Invalid params",
	CodeUserNotExists:     "User Not Exists",
	CodeUserAlreadyExists: "User Already Exists",
	CodeInvalidPassword:   "Invalid password",
	CodeServerBusy:        "Server Busy",
	CodeLoginFail:         "Login Fail",
	CodeSignUpFail:        "Sign Up Fail",
	CodeNeedAuth:          "Need Auth",
	CodeInvalidToken:      "Invalid Token",
	CodeTokenExpired:      "Token Expired",
}

func (c Rescode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
