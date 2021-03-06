package main

import (
	"flag"
	"log"
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
	port := flag.String("p", "8080", "port to serve on")
	flag.Parse()

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

	// Handlers with Auth middleware
	authMux := mux.NewRouter()
	authMux.HandleFunc("/api/v1/user/logout", userHandler.Logout).Methods("DELETE")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.GetProfile).Methods("GET")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.UpdateProfile).Methods("PUT")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.GetProfileAvatar).Methods("GET")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.UpdateProfileAvatar).Methods("PUT")
	handlersWithAuth := middleware.Auth(sessionManager, authMux)

	mainMux.Handle("/api/v1/user/", handlersWithAuth)
	mainMux.HandleFunc("/api/v1/user/signup", userHandler.Signup).Methods("POST")
	mainMux.HandleFunc("/api/v1/user/login", userHandler.Login).Methods("POST")

	// Base middlewares
	mux := middleware.Cors(mainMux)
	mux = middleware.Panic(mux)

	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
