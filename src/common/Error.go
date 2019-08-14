package common

import "errors"

var (
	LOCK_ALREADY_USED = errors.New("锁已被占用")

)