package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	session_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session/usecase"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/proto/session"
	session_server "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/server"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func InitSessionService() {
	// Load session service redis environment
	err := godotenv.Load(configs.SessionServiceRedisEnv)
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

	lis, err := net.Listen("tcp", ":8228")
	if err != nil {
		log.Fatalf("error start session service %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterSessionServiceServer(server, sessionServer)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
