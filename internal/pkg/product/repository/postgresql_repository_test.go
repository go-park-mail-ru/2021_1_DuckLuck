package repository

import (
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestPostgresqlRepository_SelectProductById(t *testing.T) {
	productId := uint64(1)
	testProduct := models.Product{
		Id:    productId,
		Title: "test",
		Price: models.ProductPrice{
			Discount: 10,
			BaseCost: 20,
		},
		Rating:       4,
		Description:  "description",
		Category:     4,
		CategoryPath: nil,
		Images:       []string{"/product/6026466446.jpg", "/product/6043224204.jpg", "/product/6043224631.jpg"},
	}

	t.Run("SelectProductById_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		productRows := sqlmock.
			NewRows([]string{"id", "title", "rating", "description", "base_cost", "discount",
				"images", "id_category"}).
			AddRow(testProduct.Id, testProduct.Title, testProduct.Rating, testProduct.Description,
				testProduct.Price.BaseCost, testProduct.Price.Discount,
				pq.Array(testProduct.Images), testProduct.Category)
		sqlMock.ExpectQuery("SELECT").WithArgs(productId).WillReturnRows(productRows)

		data, err := repository.SelectProductById(productId)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if !reflect.DeepEqual(*data, testProduct) {
			t.Errorf("not match: [%v] - [%v]", *data, testProduct)
			return
		}
	})

	t.Run("SelectProductById_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").WithArgs(productId).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectProductById(productId)
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

func TestPostgresqlRepository_SelectRangeProducts(t *testing.T) {
	categories := []uint64{1, 2, 3, 4}
	viewProduct := models.ViewProduct{
		Id:    1,
		Title: "tesst",
		Price: models.ProductPrice{
			Discount: 2,
			BaseCost: 100,
		},
		Rating:       5,
		PreviewImage: "http//test.png",
	}
	paginator := models.PaginatorProducts{
		PageNum:       2,
		Count:         10,
		SortKey:       "cost",
		SortDirection: "ASC",
		Category:      0,
	}

	// Success result
	paginators := []models.PaginatorProducts{
		models.PaginatorProducts{
			PageNum:       2,
			Count:         10,
			SortKey:       "cost",
			SortDirection: "DESC",
			Category:      0,
		},
		models.PaginatorProducts{
			PageNum:       2,
			Count:         10,
			SortKey:       "cost",
			SortDirection: "ASC",
			Category:      0,
		},
		models.PaginatorProducts{
			PageNum:       2,
			Count:         10,
			SortKey:       "rating",
			SortDirection: "ASC",
			Category:      0,
		},
		models.PaginatorProducts{
			PageNum:       2,
			Count:         10,
			SortKey:       "rating",
			SortDirection: "DESC",
			Category:      0,
		},
	}

	t.Run("SelectRangeProducts_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		for i := 0; i < 4; i++ {
			countPages := sqlmock.
				NewRows([]string{"count"}).
				AddRow(10)
			sqlMock.ExpectQuery("SELECT").
				WillReturnRows(countPages)

			category := sqlmock.
				NewRows([]string{"id", "title", "base_cost", "discount", "rating", "images[1]"}).
				AddRow(viewProduct.Id, viewProduct.Title, viewProduct.Price.BaseCost,
					viewProduct.Price.Discount, viewProduct.Rating, viewProduct.PreviewImage)
			sqlMock.ExpectQuery("SELECT").WithArgs(pq.Array(categories), paginators[i].Count,
				paginators[i].Count*(paginators[i].PageNum-1)).
				WillReturnRows(category)

			data, err := repository.SelectRangeProducts(&paginators[i], &categories)
			if err != nil {
				t.Errorf("internal err: %s", err)
				return
			}

			err = sqlMock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("expectations were not met in order: %s", err)
				return
			}

			if !reflect.DeepEqual(viewProduct, *data.ListPreviewProducts[0]) {
				t.Errorf("not match: [%v] - [%v]", viewProduct, *data.ListPreviewProducts[0])
				return
			}
		}
	})

	t.Run("SelectRangeProducts_incorrect_count_pages", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectRangeProducts(&paginator, &categories)
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

	t.Run("SelectRangeProducts_incorrect_paginator", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		countPages := sqlmock.
			NewRows([]string{"count"}).
			AddRow(10)
		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(countPages)

		sqlMock.ExpectQuery("SELECT").
			WillReturnError(sqlmock.ErrCancelled)

		_, err = repository.SelectRangeProducts(&paginator, &categories)
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

	t.Run("SelectRangeProducts_bad_options", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		countPages := sqlmock.
			NewRows([]string{"count"}).
			AddRow(10)
		sqlMock.ExpectQuery("SELECT").
			WillReturnRows(countPages)

		category := sqlmock.
			NewRows([]string{"id", "title", "base_cost"}).
			AddRow(viewProduct.Id, viewProduct.Title, viewProduct.Price.BaseCost)
		sqlMock.ExpectQuery("SELECT").WithArgs(pq.Array(categories), paginator.Count,
			paginator.Count*(paginator.PageNum-1)).
			WillReturnRows(category)

		_, err = repository.SelectRangeProducts(&paginator, &categories)
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
