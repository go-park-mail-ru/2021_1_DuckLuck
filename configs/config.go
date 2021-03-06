package configs

import "os"

var (
	PathProject, _         = os.Getwd()
	PathToUploadAvatar     = PathProject + "/uploads/avatar"
	PathToUploadProductImg = PathProject + "/uploads/product"
	UrlToAvatar            = PathProject + "/avatar/"
	UrlToProductImg        = PathProject + "/product/"

	CorsOrigins = map[string]struct{}{
		"http://localhost":               struct{}{},
		"http://localhost:3000":          struct{}{},
		"http://duckluckmarket.xyz":      struct{}{},
		"http://duckluckmarket.xyz:3000": struct{}{},
	}
)
