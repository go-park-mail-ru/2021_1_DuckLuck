package user

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user Handler

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfileAvatar(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	GetProfileAvatar(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}
