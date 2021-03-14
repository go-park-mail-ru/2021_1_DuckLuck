package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

const goodSessionValue = "goodValue"

var retSession = &models.Session{
	Value: goodSessionValue,
}

var err error

func TestLocalRepository_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().SelectSessionByValue(goodSessionValue).Times(1).Return(retSession, nil)

	useCase := NewUseCase(mockRepository)
	var res *models.Session

	res, err = useCase.Check(goodSessionValue)
	require.NoError(t, err)
	require.Equal(t, res, retSession)
}
