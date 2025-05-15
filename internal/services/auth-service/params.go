package authservice

import (
	"time"
)

var (
	SessionExpireTime = time.Hour
)

type authContextKey int

const (
	AuthContextKey    authContextKey = iota
	sessionContextKey authContextKey = iota
)

func UpdateSessionExpireTime(t time.Duration) {
	if t <= 0 {
		panic("session expire time must be greater than 0")
	}
	SessionExpireTime = t
}
