package models

const (
	SessionCookieName   = "session_id"
	SessionContextKey   = "session_key"
	RequireIdKey        = "require_key"
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
