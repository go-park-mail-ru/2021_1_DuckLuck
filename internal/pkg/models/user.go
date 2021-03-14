package models

type ProfileUser struct {
	Id        uint64 `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Avatar    Avatar `json:"avatar"`
}

type Avatar struct {
	Url		string `json:"url"`
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
