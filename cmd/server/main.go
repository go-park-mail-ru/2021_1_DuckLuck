package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	_ "github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	cart_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/handler"
	cart_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/repository"
	cart_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/usecase"
	category_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/handler"
	category_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/repository"
	category_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/usecase"
	csrf_token_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token/handler"
	order_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/handler"
	order_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/repository"
	order_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/usecase"
	product_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/handler"
	product_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/repository"
	product_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/usecase"
	session_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/handler"
	session_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_delivery "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/handler"
	user_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/repository"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/middleware"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/s3_utils"
	_ "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/s3_utils"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitApiServer() {
	// Load server api environment
	err := godotenv.Load(configs.PathToApiServerEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load postgresql environment
	err = godotenv.Load(configs.PathToPostgreSqlEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load redis environment
	err = godotenv.Load(configs.PathToRedisEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	mainLogger := logger.Logger{}
	err = mainLogger.InitLogger()
	if err != nil {
		log.Fatal(errors.ErrOpenFile.Error())
	}
}

func main() {
	InitApiServer()
	s3_utils.InitS3()

	// Connect to postgreSql db
	postgreSqlConn, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=%s "+
			"port=%s "+
			"sslmode=%s ",
			os.Getenv("PG_USER"),
			os.Getenv("PG_PASS"),
			os.Getenv("PG_DB_NAME"),
			os.Getenv("PG_HOST"),
			os.Getenv("PG_PORT"),
			os.Getenv("PG_SSL_MODE"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()
	if err := postgreSqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	// Connect to redis db
	redisConn := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if redisConn == nil {
		log.Fatal(errors.ErrDBFailedConnection.Error())
	}
	defer redisConn.Close()

	sessionRepo := session_repo.NewSessionRedisRepository(redisConn)
	sessionUCase := session_usecase.NewUseCase(sessionRepo)
	sessionHandler := session_delivery.NewHandler(sessionUCase)

	categoryRepo := category_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	categoryUCase := category_usecase.NewUseCase(categoryRepo)
	categoryHandler := category_delivery.NewHandler(categoryUCase)

	productRepo := product_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	productUCase := product_usecase.NewUseCase(productRepo, categoryRepo)
	productHandler := product_delivery.NewHandler(productUCase)

	cartRepo := cart_repo.NewSessionRedisRepository(redisConn)
	cartUCase := cart_usecase.NewUseCase(cartRepo, productRepo)
	cartHandler := cart_delivery.NewHandler(cartUCase)

	userRepo := user_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	userUCase := user_usecase.NewUseCase(userRepo)
	userHandler := user_delivery.NewHandler(userUCase, sessionUCase)

	orderRepo := order_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	orderUCase := order_usecase.NewUseCase(orderRepo, cartRepo, productRepo, userRepo)
	orderHandler := order_delivery.NewHandler(orderUCase, cartUCase)

	csrfTokenHandler := csrf_token_delivery.NewHandler()

	mainMux := mux.NewRouter()
	mainMux.Use(middleware.AccessLog)
	mainMux.Use(middleware.Panic)
	mainMux.Use(middleware.Cors)
	// Check csrf token
	mainMux.Use(middleware.CsrfCheck)

	mainMux.HandleFunc("/api/v1/csrf", csrfTokenHandler.GetCsrfToken).Methods("GET", "OPTIONS")
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
	authMux.HandleFunc("/api/v1/session", sessionHandler.CheckSession).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/order", orderHandler.GetUserOrders).Methods("POST", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/logout", userHandler.Logout).Methods("DELETE", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile", userHandler.UpdateProfile).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.GetProfileAvatar).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/user/profile/avatar", userHandler.UpdateProfileAvatar).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart", cartHandler.GetProductsFromCart).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart", cartHandler.DeleteProductsFromCart).Methods("DELETE", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.ChangeProductInCart).Methods("PUT", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.AddProductInCart).Methods("POST", "OPTIONS")
	authMux.HandleFunc("/api/v1/cart/product", cartHandler.DeleteProductInCart).Methods("DELETE", "OPTIONS")
	authMux.HandleFunc("/api/v1/order", orderHandler.GetOrderFromCart).Methods("GET", "OPTIONS")
	authMux.HandleFunc("/api/v1/order", orderHandler.AddCompletedOrder).Methods("POST", "OPTIONS")

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("API_SERVER_HOST"),
			os.Getenv("API_SERVER_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mainMux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
