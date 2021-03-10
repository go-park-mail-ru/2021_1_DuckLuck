package repository

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

const sessionValue = "value"

var sessionForAdd = &models.Session{
	Value:  "value",
	UserId: 0,
}

var err error
var rep = NewSessionLocalRepository()

func TestLocalRepository_Add(t *testing.T) {
	err = rep.Add(sessionForAdd)
	require.NoError(t, err)
}

func TestLocalRepository_GetByValue(t *testing.T) {
	_, err = rep.GetByValue(sessionValue)
	require.NoError(t, err)
}

func TestLocalRepository_DestroyByValue(t *testing.T) {
	_, err = rep.GetByValue(sessionValue)
	require.NoError(t, err)
}
