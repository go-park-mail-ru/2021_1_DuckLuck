package main

import (
	"fmt"

	"log"
	"net"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/metrics"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/grpc_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
	session_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session/usecase"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/proto/session"
	session_server "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/server"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func InitSessionService() {
	// Load session_service environment
	err := godotenv.Load(configs.SessionServiceMainEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load session service redis environment
	err = godotenv.Load(configs.SessionServiceRedisEnv)
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
	err = mainLogger.InitLogger(configs.SessionServiceLog, os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitSessionService()

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
	sessionServer := session_server.NewSessionServer(sessionUCase)

	lis, err := net.Listen(
		os.Getenv("SESSION_SERVICE_PROTOCOL"),
		fmt.Sprintf("%s:%s",
			os.Getenv("SESSION_SERVICE_HOST"),
			os.Getenv("SESSION_SERVICE_PORT")),
	)
	if err != nil {
		log.Fatalf("error start session service %v", err)
	}

	metric, err := metrics.CreateNewMetrics("session_service")
	if err != nil {
		log.Fatal(err)
	}
	accessInterceptor := grpc_utils.AccessInterceptor(metric)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor),
	)
	proto.RegisterSessionServiceServer(server, sessionServer)

	go metrics.CreateNewMetricsRouter(os.Getenv("SESSION_SERVICE_HOST"))

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
