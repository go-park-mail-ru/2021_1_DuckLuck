package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/metrics"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/grpc_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
	auth_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user/repository"
	auth_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user/usecase"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/proto/user"
	auth_server "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func InitAuthService() {
	// Load auth_service environment
	err := godotenv.Load(configs.AuthServiceMainEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load auth service postgresql environment
	err = godotenv.Load(configs.AuthServicePostgreSqlEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load network environment
	err = godotenv.Load(configs.NetworkEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	mainLogger := logger.Logger{}
	err = mainLogger.InitLogger(configs.AuthServiceLog, os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitAuthService()

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

	authRepo := auth_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	authUCase := auth_usecase.NewUseCase(authRepo)
	authServer := auth_server.NewAuthServer(authUCase)

	lis, err := net.Listen(
		os.Getenv("AUTH_SERVICE_PROTOCOL"),
		fmt.Sprintf("%s:%s",
			os.Getenv("AUTH_SERVICE_HOST"),
			os.Getenv("AUTH_SERVICE_PORT")),
	)
	if err != nil {
		log.Fatalf("error start session service %v", err)
	}

	metric, err := metrics.CreateNewMetrics("auth_service")
	if err != nil {
		log.Fatal(err)
	}
	accessInterceptor := grpc_utils.AccessInterceptor(metric)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor),
	)
	proto.RegisterAuthServiceServer(server, authServer)

	go metrics.CreateNewMetricsRouter(os.Getenv("AUTH_SERVICE_HOST"))

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
