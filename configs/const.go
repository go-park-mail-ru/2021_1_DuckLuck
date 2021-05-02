package configs

import (
	"os"
)

var (
	PathProject, _ = os.Getwd()

	NetworkEnv = PathProject + "/configs/envs/network.env"

	ApiServerMainEnv       = PathProject + "/configs/envs/api_server/main.env"
	ApiServerPostgreSqlEnv = PathProject + "/configs/envs/api_server/postgresql.env"
	ApiServerRedisEnv      = PathProject + "/configs/envs/api_server/redis.env"
	ApiServerS3Env         = PathProject + "/configs/envs/api_server/s3.env"
	ApiServerLog           = PathProject + "/log/api_server_log.txt"

	SessionServiceMainEnv  = PathProject + "/configs/envs/session_service/main.env"
	SessionServiceRedisEnv = PathProject + "/configs/envs/session_service/redis.env"
	SessionServiceLog      = PathProject + "/log/session_service_log.txt"

	CartServiceMainEnv  = PathProject + "/configs/envs/cart_service/main.env"
	CartServiceRedisEnv = PathProject + "/configs/envs/cart_service/redis.env"
	CartServiceLog      = PathProject + "/log/cart_service_log.txt"

	AuthServiceMainEnv       = PathProject + "/configs/envs/auth_service/main.env"
	AuthServicePostgreSqlEnv = PathProject + "/configs/envs/auth_service/postgresql.env"
	AuthServiceLog           = PathProject + "/log/auth_service_log.txt"

	CorsOrigin = "http://localhost:3000"
)
