package models

type ProfileUser struct {
	Id        uint64 `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Avatar    string `json:"-"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
