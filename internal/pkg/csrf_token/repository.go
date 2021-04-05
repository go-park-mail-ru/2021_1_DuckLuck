package csrf_token

type Repository interface {
	AddCsrfToken(tokenValue string) error
	CheckCsrfToken(tokenValue string) bool
}
