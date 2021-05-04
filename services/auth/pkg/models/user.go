package models

// Model contains fields for login
type LoginUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

// Model contains fields for registration new user
type SignupUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

// All data for user authorize
type AuthUser struct {
	Id       uint64 `json:"-"`
	Email    string `json:"email" valid:"email"`
	Password []byte `json:"-"`
}
