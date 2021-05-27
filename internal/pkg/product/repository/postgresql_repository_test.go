package repository

import (
	"github.com/stretchr/testify/assert"
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
			Discount:  10,
			TotalCost: 18,
			BaseCost:  20,
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
			NewRows([]string{"id", "title", "description",
				"properties", "base_cost", "total_cost", "discount",
				"images", "id_category", "avg_rating", "count_reviews"}).
			AddRow(testProduct.Id, testProduct.Title, testProduct.Description,
				testProduct.Properties, testProduct.Price.BaseCost, testProduct.Price.TotalCost,
				testProduct.Price.Discount, pq.Array(testProduct.Images), testProduct.Category,
				testProduct.Rating, testProduct.CountReviews)
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
	userId := uint64(3)
	product := models.ViewProduct{}
	paginator := models.PaginatorProducts{}
	sortString := ""
	filterString := ""

	t.Run("SelectRangeProducts_success", func(t *testing.T) {
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
			WithArgs(paginator.Category, paginator.Count, paginator.Count*(paginator.PageNum-1)).
			WillReturnRows(productRows)

		favorites, err := repository.SelectRangeProducts(&paginator, sortString, filterString)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, []*models.ViewProduct{&product}, favorites)
	})

	t.Run("SelectRangeProducts_internal_db_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.ExpectQuery("SELECT").
			WithArgs(userId, paginator.Count, paginator.Count*(paginator.PageNum-1)).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.SelectRangeProducts(&paginator, sortString, filterString)
		assert.Error(t, err, "expected error")
	})
}
