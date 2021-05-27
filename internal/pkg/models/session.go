package models

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	SessionCookieName   = "session_id"
	SessionContextKey   = contextKey("session_key")
	RequireIdKey        = contextKey("require_key")
	RequireIdName       = "require_id"
	ExpireSessionCookie = 90 * 24 * 3600
)

type Session struct {
	Value    string
	UserData UserId
}

type UserId struct {
	Id uint64
}
