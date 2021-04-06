package configs

import "os"

var (
	PathProject, _ = os.Getwd()

	PathToUpload    = PathProject + "/uploads"
	UrlToAvatar     = "/avatar/"
	UrlToProductImg = "/product/"

	PathToLogFile = PathProject + "/log/log.txt"
	LogLevel      = "debug"

	CorsOrigin = "https://duckluckmarket.xyz"
)
