package configs

import (
	"os"
)

var (
	PathProject, _ = os.Getwd()

	PathToLogFile = PathProject + "/log/log.txt"
	LogLevel      = "debug"

	PathToApiServerEnv  = PathProject + "/configs/envs/main_server/api_server.env"
	PathToPostgreSqlEnv = PathProject + "/configs/envs/main_server/postgresql.env"
	PathToRedisEnv      = PathProject + "/configs/envs/main_server/redis.env"
	PathToS3Env         = PathProject + "/configs/envs/main_server/s3.env"

	SessionServiceRedisEnv = PathProject + "/configs/envs/session_service/redis.env"

	CorsOrigin = "http://localhost:3000"
)
