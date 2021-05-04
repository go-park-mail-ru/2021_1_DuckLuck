package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlRepository_SelectProfileById(t *testing.T) {
	userId := uint64(3)
	unknownUserId := uint64(1)
	userProfile := models.ProfileUser{
		Id:        userId,
		FirstName: "name",
		LastName:  "last_name",
		Avatar: models.Avatar{
			Url: "http://test.png",
		},
	}

	t.Run("SelectProfileById_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "first_name", "last_name", "email", "password", "avatar"}).
			AddRow(userProfile.Id, userProfile.FirstName, userProfile.LastName, userProfile.Avatar.Url)
		sqlMock.ExpectQuery("SELECT").WithArgs(userId).WillReturnRows(rows)

		data, err := repository.SelectProfileById(userId)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if !reflect.DeepEqual(*data, userProfile) {
			t.Errorf("not match: [%v] - [%v]", *data, userProfile)
			return
		}
	})

	t.Run("SelectProfileById_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(userId).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectProfileById(userId)
		if err != errors.ErrDBInternalError {
			t.Error("expected error")
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})

	t.Run("SelectProfileById_user_not_found", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(unknownUserId).WillReturnError(sql.ErrNoRows)

		_, err = repository.SelectProfileById(unknownUserId)
		if err != errors.ErrUserNotFound {
			t.Error("expected error")
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})
}

func TestPostgresqlRepository_UpdateProfile(t *testing.T) {
	userId := uint64(2)

	updateInfo := models.UpdateUser{
		FirstName: "new_name",
		LastName:  "new_last_name",
	}

	t.Run("UpdateProfile_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectExec("UPDATE users").WithArgs(updateInfo.FirstName,
			updateInfo.LastName, userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repository.UpdateProfile(userId, &updateInfo)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})

	t.Run("UpdateProfile_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectExec("UPDATE users").WithArgs(updateInfo.FirstName,
			updateInfo.LastName, userId).
			WillReturnError(sql.ErrConnDone)

		err = repository.UpdateProfile(userId, &updateInfo)
		if err == nil {
			t.Error("expected error")
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})
}

func TestPostgresqlRepository_UpdateAvatar(t *testing.T) {

	userId := uint64(2)

	avatarUrl := "https://test.png"

	t.Run("UpdateAvatar_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectExec("UPDATE users").WithArgs(avatarUrl, userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repository.UpdateAvatar(userId, avatarUrl)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})

	t.Run("UpdateAvatar_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectExec("UPDATE users").WithArgs(avatarUrl, userId).
			WillReturnError(sql.ErrConnDone)

		err = repository.UpdateAvatar(userId, avatarUrl)
		if err == nil {
			t.Error("expected error")
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})
}
