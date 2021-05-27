package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresqlRepository_AddProductToFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)

	t.Run("AddProductToFavorites_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		sqlMock.ExpectExec("INSERT INTO favorites").
			WithArgs(productId, userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repository.AddProductToFavorites(productId, userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddProductToFavorites_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		sqlMock.ExpectExec("INSERT INTO favorites").
			WithArgs(productId, userId).
			WillReturnError(errors.ErrDBInternalError)

		err = repository.AddProductToFavorites(productId, userId)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_DeleteProductFromFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)

	t.Run("DeleteProductFromFavorites_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		sqlMock.ExpectExec("DELETE").
			WithArgs(productId, userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repository.DeleteProductFromFavorites(productId, userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("DeleteProductFromFavorites_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		sqlMock.ExpectExec("DELETE").
			WithArgs(productId, userId).
			WillReturnError(errors.ErrDBInternalError)

		err = repository.DeleteProductFromFavorites(productId, userId)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_GetCountPages(t *testing.T) {
	userId := uint64(3)
	span := 10
	countFavorites := 3

	t.Run("GetCountPages_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		productRows := sqlmock.
			NewRows([]string{"count"}).
			AddRow(countFavorites)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnRows(productRows)

		countPages, err := repository.GetCountPages(userId, span)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, 1, countPages)
	})

	t.Run("GetCountPages_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.GetCountPages(userId, span)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_CreateSortString(t *testing.T) {
	discountKey := "discount"
	dateKey := "date"
	ratingKey := "rating"
	costKey := "cost"
	badSortKey := "price"
	DescDirection := "DESC"
	AscDirection := "ASC"
	badSortDirection := "test"

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		_, err = repository.CreateSortString(discountKey, DescDirection)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		_, err = repository.CreateSortString(dateKey, AscDirection)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		_, err = repository.CreateSortString(ratingKey, DescDirection)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		_, err = repository.CreateSortString(costKey, DescDirection)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("CreateSortString_internal_db_error", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		_, err = repository.CreateSortString(badSortKey, badSortDirection)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_SelectRangeFavorites(t *testing.T) {
	userId := uint64(3)
	product := models.ViewFavorite{}
	paginator := models.PaginatorFavorites{}
	sortString := ""

	t.Run("SelectRangeFavorites_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		productRows := sqlmock.
			NewRows([]string{"id", "title", "base_cost", "total_cost",
				"discount", "images", "avg_rating", "count_reviews"}).
			AddRow(product.Id, product.Title, product.Price.BaseCost,
				product.Price.TotalCost, product.Price.Discount,
				product.PreviewImage, product.Rating, product.CountReviews)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId, paginator.Count, paginator.Count*(paginator.PageNum-1)).
			WillReturnRows(productRows)

		favorites, err := repository.SelectRangeFavorites(&paginator, sortString, userId)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, []*models.ViewFavorite{&product}, favorites)
	})

	t.Run("SelectRangeFavorites_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId, paginator.Count, paginator.Count*(paginator.PageNum-1)).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectRangeFavorites(&paginator, sortString, userId)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_GetUserFavorites(t *testing.T) {
	userId := uint64(3)
	productId := uint64(4)

	t.Run("GetUserFavorites_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		productRows := sqlmock.
			NewRows([]string{"product_id"}).
			AddRow(productId)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnRows(productRows)

		_, err = repository.GetUserFavorites(userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetUserFavorites_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.GetUserFavorites(userId)
		assert.Error(t, err, "unexpected error")
	})
}
