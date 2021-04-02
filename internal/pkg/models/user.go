package models

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
)

// All user information
// This models saved in database
type ProfileUser struct {
	Id        uint64         `json:"-"`
	FirstName sql.NullString `json:"first_name" valid:"utfletter, stringlength(3|30)"`
	LastName  sql.NullString `json:"last_name" valid:"utfletter, stringlength(3|30)"`
	Email     string         `json:"email" valid:"email"`
	Password  []byte         `json:"-"`
	Avatar    Avatar         `json:"avatar" valid:"notnull, json"`
}

// User avatar
type Avatar struct {
	Url sql.NullString `json:"url" valid:"minstringlength(3)"`
}

// Model contains fields for updating user information
type UpdateUser struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(3|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(3|30)"`
}

func (uu *UpdateUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	uu.FirstName = sanitizer.Sanitize(uu.FirstName)
	uu.LastName = sanitizer.Sanitize(uu.LastName)
}

// Model contains fields for login
type LoginUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

func (lu *LoginUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	lu.Email = sanitizer.Sanitize(lu.Email)
	lu.Password = sanitizer.Sanitize(lu.Password)
}

// Model contains fields for registration new user
type SignupUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

func (su *SignupUser) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	su.Email = sanitizer.Sanitize(su.Email)
	su.Password = sanitizer.Sanitize(su.Password)
}
