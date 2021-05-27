package admin

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/admin UseCase

type UseCase interface {
	ChangeOrderStatus(updateOrder *models.UpdateOrder) error
}
