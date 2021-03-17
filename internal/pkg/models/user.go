package models

import "database/sql"

type ProfileUser struct {
	Id        uint64         `json:"-"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	Avatar    Avatar         `json:"avatar"`
}

type Avatar struct {
	Url sql.NullString `json:"url"`
}

type UpdateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
