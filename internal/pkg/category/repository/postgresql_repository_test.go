package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlRepository_GetNextLevelCategories(t *testing.T) {
	categoryId := uint64(3)
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}

	t.Run("GetNextLevelCategories_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "name"}).
			AddRow(categories[0].Id, categories[0].Name)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)

		data, err := repository.GetNextLevelCategories(categoryId)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if !reflect.DeepEqual(data, categories) {
			t.Errorf("not match: [%v] - [%v]", data, categories)
			return
		}
	})

	t.Run("GetNextLevelCategories_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnError(sql.ErrConnDone)

		_, err = repository.GetNextLevelCategories(categoryId)
		if err != errors.ErrIncorrectPaginator {
			t.Error("expected error")
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}
	})

	t.Run("GetNextLevelCategories_can't_scan_rows", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(categories[0].Id)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)
		_, err = repository.GetNextLevelCategories(categoryId)
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

func TestPostgresqlRepository_GetCategoriesByLevel(t *testing.T) {
	categoryId := uint64(3)
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}

	t.Run("GetCategoriesByLevel_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "name"}).
			AddRow(categories[0].Id, categories[0].Name)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)

		data, err := repository.GetCategoriesByLevel(categoryId)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if !reflect.DeepEqual(data, categories) {
			t.Errorf("not match: [%v] - [%v]", data, categories)
			return
		}
	})

	t.Run("GetCategoriesByLevel_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnError(sql.ErrConnDone)

		_, err = repository.GetCategoriesByLevel(categoryId)
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

	t.Run("GetCategoriesByLevel_can't_scan_rows", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(categories[0].Id)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)
		_, err = repository.GetCategoriesByLevel(categoryId)
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

func TestPostgresqlRepository_GetAllSubCategoriesId(t *testing.T) {

}

func TestPostgresqlRepository_GetPathToCategory(t *testing.T) {
	categoryId := uint64(3)
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}

	t.Run("GetPathToCategory_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "name"}).
			AddRow(categories[0].Id, categories[0].Name)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)

		data, err := repository.GetPathToCategory(categoryId)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if !reflect.DeepEqual(data, categories) {
			t.Errorf("not match: [%v] - [%v]", data, categories)
			return
		}
	})

	t.Run("GetPathToCategory_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnError(sql.ErrConnDone)

		_, err = repository.GetPathToCategory(categoryId)
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

	t.Run("GetPathToCategory_can't_scan_rows", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow(categories[0].Id)
		sqlMock.ExpectQuery("SELECT").WithArgs(categoryId).
			WillReturnRows(rows)
		_, err = repository.GetPathToCategory(categoryId)
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
