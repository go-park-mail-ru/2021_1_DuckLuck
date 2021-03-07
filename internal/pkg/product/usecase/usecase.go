package product

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
)

type ProductUseCase struct {
}

func NewUseCase() product.UseCase {
	return &ProductUseCase{}
}
