package configs

var (
	PathToUploads = "/home/uploads"
	CorsOrigins   = map[string]struct{}{
		"http://localhost": struct{}{},
		"http://localhost:3000": struct{}{},
		"http://duckluckmarket.xyz": struct{}{},
		"http://duckluckmarket.xyz:3000": struct{}{},
	}
)
