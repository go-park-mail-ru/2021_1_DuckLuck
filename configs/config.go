package configs

import "os"

var (
	PathProject, _ = os.Getwd()

	PathToUpload    = PathProject + "/uploads"
	UrlToAvatar     = "/avatar/"
	UrlToProductImg = "/product/"

	PathToLogFile = PathProject + "/log/log.txt"
	LogLevel      = "debug"

	PathToApiServerEnv  = PathProject + "/configs/envs/api_server.env"
	PathToPostgreSqlEnv = PathProject + "/configs/envs/postgresql.env"
	PathToRedisEnv      = PathProject + "/configs/envs/redis.env"
	PathToS3Env         = PathProject + "/configs/envs/s3.env"

	CorsOrigin = "http://localhost:3000"
)
