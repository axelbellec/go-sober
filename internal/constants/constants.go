package constants

import "time"

type contextKey string

const (
	// UserContextKey is used to store/retrieve
	// user information in the request context
	UserContextKey contextKey = "go-sober-user"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 50
)

var (
	DefaultStartDate = time.Unix(0, 0)
	DefaultEndDate   = time.Now()
)
