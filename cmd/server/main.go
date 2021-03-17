package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	product_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/handler"
	product_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/repository"
	product_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/usecase"
	session_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/handler"
	user_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/repository"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/middleware"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	port := flag.String("p", "8080", "port to serve on")
	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	flag.Parse()

	// Database
	pgConn, err := sql.Open(
		"postgres",
		"user=postgres "+
		"password=Id47806649 "+
		"dbname=ozon_db "+
		"host=localhost "+
		"port=5432",
	)

	if err != nil {
		log.Fatal(err)
	}
	defer pgConn.Close()

	if err := pgConn.Ping(); err != nil {
		log.Fatal(err)
	}

	c, err := redis.DialURL(*redisAddr)
	if err != nil {
		panic(errors.ErrDBFailedConnection.Error())
	}
	defer c.Close()

	sessionRepo := session_repo.NewSessionRedisRepository(c)
	sessionUCase := session_usecase.NewUseCase(sessionRepo)

	userRepo := user_repo.NewSessionPostgresqlRepository(pgConn)
	userUCase := user_usecase.NewUseCase(userRepo)
	userHandler := user_delivery.NewHandler(userUCase, sessionUCase)

	productRepo := product_repo.NewSessionPostgresqlRepository(pgConn)
	productUCase := product_usecase.NewUseCase(productRepo)
	productHandler := product_delivery.NewHandler(productUCase)

	mainMux := mux.NewRouter()
	mainMux.Use(middleware.Panic)
	mainMux.Use(middleware.Cors)
	mainMux.HandleFunc("/api/v1/user/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	mainMux.HandleFunc("/api/v1/user/login", userHandler.Login).Methods("POST", "OPTIONS")
	mainMux.HandleFunc("/api/v1/product/{id:[0-9]+}", productHandler.GetProduct).Methods("GET", "OPTIONS")
	mainMux.HandleFunc("/api/v1/product", productHandler.GetListPreviewProducts).Methods("POST", "OPTIONS")

	// Handlers with Auth middleware
	authMux := mainMux.PathPrefix("/").Subrouter()
	middlewareAuth := middleware.Auth(sessionUCase)
	authMux.Use(middlewareAuth)
	authMux.HandleFunc("/api/v1/user/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/logout", userHandler.Logout).Methods("DELETE", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.UpdateProfile).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.GetProfileAvatar).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.UpdateProfileAvatar).Methods("PUT", "OPTIONS")

	server := &http.Server{
		Addr:         ":" + *port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mainMux,
	}

	fmt.Println("starting server")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
