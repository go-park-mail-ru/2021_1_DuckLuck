package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	cart_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/cart/repository"
	cart_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/cart/usecase"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/proto/cart"
	cart_server "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/server"
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

	cartRepo := cart_repo.NewSessionRedisRepository(redisConn)
	cartUCase := cart_usecase.NewUseCase(cartRepo)
	cartServer := cart_server.NewCartServer(cartUCase)

	lis, err := net.Listen("tcp", ":8448")
	if err != nil {
		log.Fatalf("error start session service %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterCartServiceServer(server, cartServer)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
