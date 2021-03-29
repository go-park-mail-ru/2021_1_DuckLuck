package models

import "database/sql"

type ProfileUser struct {
	Id        uint64         `json:"-"`
	FirstName sql.NullString `json:"first_name" valid:"utfletter, stringlength(3|30)"`
	LastName  sql.NullString `json:"last_name" valid:"utfletter, stringlength(3|30)"`
	Email     string         `json:"email" valid:"email"`
  Password  []byte         `json:"-"`
	Avatar    Avatar         `json:"avatar" valid:"notnull, json"`
}

type Avatar struct {
	Url sql.NullString `json:"url" valid:"minstringlength(3)"`
}

type UpdateUser struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(3|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(3|30)"`
}

type LoginUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}

type SignupUser struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"stringlength(6|32)"`
}
