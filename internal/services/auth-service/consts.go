package authservice

import (
	"time"
)

const (
	SessionExpireTime = time.Hour
)

type authContextKey int

const (
	AuthContextKey authContextKey = iota
)
