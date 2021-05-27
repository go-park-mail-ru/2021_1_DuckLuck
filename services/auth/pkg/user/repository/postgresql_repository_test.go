package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/tools/password_hasher"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlRepository_SelectUserByEmail(t *testing.T) {
	userId := uint64(3)

	userEmail := "test@test.ru"
	unknownUserEmail := "bad@email.ru"

	hashedPassword, _ := password_hasher.GenerateHashFromPassword("qwerty")

	userProfile := models.AuthUser{
		Id:       userId,
		Password: hashedPassword,
		Email:    userEmail,
	}

	t.Run("SelectUserByEmail_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "email", "password"}).
			AddRow(userProfile.Id, userProfile.Email, userProfile.Password)
		sqlMock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)

		data, err := repository.SelectUserByEmail(userEmail)
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

	t.Run("SelectUserByEmail_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(userEmail).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectUserByEmail(userEmail)
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

	t.Run("SelectUserByEmail_user_not_found", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(unknownUserEmail).WillReturnError(sql.ErrNoRows)

		_, err = repository.SelectUserByEmail(unknownUserEmail)
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

func TestPostgresqlRepository_AddProfile(t *testing.T) {
	userId := uint64(3)

	userEmail := "test@test.ru"
	//unknownUserEmail := "bad@email.ru"

	hashedPassword, _ := password_hasher.GenerateHashFromPassword("qwerty")

	userProfile := models.AuthUser{
		Id:       userId,
		Password: hashedPassword,
		Email:    userEmail,
	}

	t.Run("AddProfile_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("INSERT INTO").WithArgs(userProfile.Email, userProfile.Password).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.AddProfile(&userProfile)
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

}
