package main

import (
	"net/http"

	session_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/delivery"
	user_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/repository"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/middleware"

	"github.com/gorilla/mux"
)

func main() {
	sessionRepo := session_repo.NewSessionLocalRepository()
	sessionManager := &session_usecase.UseCase{
		SessionRepo: sessionRepo,
	}

	userRepo := user_repo.NewSessionLocalRepository()
	userUCase := user_usecase.NewUseCase(userRepo)
	userHandler := &user_delivery.UserHandler{
		UserUCase:      *userUCase,
		SessionManager: *sessionManager,
	}

	mainMux := mux.NewRouter()

	// Logout handler with Auth middleware
	logoutMux := mux.NewRouter()
	logoutHandler := middleware.Auth(sessionManager, logoutMux)
	mainMux.Handle("/api/v1/user/logout", logoutHandler).Methods("DELETE")

	mainMux.HandleFunc("/api/v1/user/signup", userHandler.Signup).Methods("POST")
	mainMux.HandleFunc("/api/v1/user/login", userHandler.Login).Methods("POST")

	// Base middlewares
	mux := middleware.Cors(mainMux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	http.ListenAndServe(addr, mux)
}
