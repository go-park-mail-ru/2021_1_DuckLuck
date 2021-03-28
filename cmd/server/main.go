package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	cart_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/handler"
	cart_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/repository"
	cart_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/usecase"
	category_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/handler"
	category_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/repository"
	category_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/usecase"
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
		"user=ozon_root "+
			"password=qwerty123 "+
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

	productRepo := product_repo.NewSessionPostgresqlRepository(pgConn)
	productUCase := product_usecase.NewUseCase(productRepo)
	productHandler := product_delivery.NewHandler(productUCase)

	cartRepo := cart_repo.NewSessionRedisRepository(c)
	cartUCase := cart_usecase.NewUseCase(cartRepo, productRepo)
	cartHandler := cart_delivery.NewHandler(cartUCase)

	userRepo := user_repo.NewSessionPostgresqlRepository(pgConn)
	userUCase := user_usecase.NewUseCase(userRepo)
	userHandler := user_delivery.NewHandler(userUCase, sessionUCase)

	categoryRepo := category_repo.NewSessionPostgresqlRepository(pgConn)
	categoryUCase := category_usecase.NewUseCase(categoryRepo)
	categoryHandler := category_delivery.NewHandler(categoryUCase)

	mainMux := mux.NewRouter()
	mainMux.Use(middleware.Panic)
	mainMux.Use(middleware.Cors)
	mainMux.HandleFunc("/api/v1/user/signup", userHandler.Signup).Methods("POST", "OPTIONS")
	mainMux.HandleFunc("/api/v1/user/login", userHandler.Login).Methods("POST", "OPTIONS")
	mainMux.HandleFunc("/api/v1/product/{id:[0-9]+}", productHandler.GetProduct).Methods("GET", "OPTIONS")
	mainMux.HandleFunc("/api/v1/product", productHandler.GetListPreviewProducts).Methods("POST", "OPTIONS")
	mainMux.HandleFunc("/api/v1/category", categoryHandler.GetCatalogCategories).Methods("GET", "OPTIONS")
	mainMux.HandleFunc("/api/v1/category/{id:[0-9]+}", categoryHandler.GetSubCategories).Methods("GET", "OPTIONS")

	// Handlers with Auth middleware
	authMux := mainMux.PathPrefix("/").Subrouter()
	middlewareAuth := middleware.Auth(sessionUCase)
	authMux.Use(middlewareAuth)
	authMux.HandleFunc("/api/v1/user/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/logout", userHandler.Logout).Methods("DELETE", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.UpdateProfile).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.GetProfileAvatar).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.UpdateProfileAvatar).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart", cartHandler.GetProductsFromCart).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.ChangeProductInCart).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.AddProductInCart).Methods("POST", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.DeleteProductInCart).Methods("DELETE", "OPTIONS")

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
