package repository

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

const sessionValue = "value"
const sessionUserId = 0

var sessionForAdd = &models.Session{
	Value:  sessionValue,
	UserId: sessionUserId,
}

var err error
var retSession *models.Session

func TestLocalRepository_Add(t *testing.T) {
	rep := NewSessionLocalRepository()

	err = rep.Add(sessionForAdd)
	require.NoError(t, err)
}

func TestLocalRepository_GetByValue(t *testing.T) {
	rep := NewSessionLocalRepository()
	err = rep.Add(sessionForAdd)

	retSession, err = rep.GetByValue(sessionValue)
	require.NoError(t, err)
	require.Equal(t, retSession, sessionForAdd)
}

func TestLocalRepository_DestroyByValue(t *testing.T) {
	rep := NewSessionLocalRepository()
	err = rep.Add(sessionForAdd)

	err = rep.DestroyByValue(sessionValue)
	require.NoError(t, err)
}
