package product

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product UseCase

type UseCase interface {
}
