package errorcode

import "errors"

var (
	LoginAuthFailed  = errors.New("account or password not correct")
	GetSessionFailed = errors.New("can not find session")
)
