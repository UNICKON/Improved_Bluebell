package errorcode

import "errors"

var (
	ErrorUserExist       = errors.New("user already exist")
	ErrorUserNotExist    = errors.New("user does not exist")
	ErrorWrongPassword   = errors.New("wrong password")
	ErrorInvalidToken    = errors.New("invalid token")
	ErrorUserNotLogin    = errors.New("user not login")
	ErrorVoteTimeExpired = errors.New("vote time expired")
	ErrorInvalidID       = errors.New("invalid id")
	ErrorServerBusy      = errors.New("server busy, please retry")
	ErrorAlreadyVote     = errors.New("already vote")
)
