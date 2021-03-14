package configs

import "os"

var (
	PathProject, _  = os.Getwd()
	PathToUpload    = PathProject + "/uploads"
	UrlToAvatar     = "/avatar/"
	UrlToProductImg = "/product/"

	CorsOrigin = "http://localhost:3000"
)
