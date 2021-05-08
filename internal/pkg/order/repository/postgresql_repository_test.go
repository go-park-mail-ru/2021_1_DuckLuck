package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresqlRepository_SelectRangeOrders(t *testing.T) {
	orders := []*models.PlacedOrder{
		{
			Id:             uint64(1),
			Address:        models.OrderAddress{
				Address: "test street",
			},
			TotalCost:      12,
			Products:       nil,
			DateAdded:      time.Time{},
			DateDelivery:   time.Time{},
			OrderNumber:    models.OrderNumber{},
			StatusPay:      "ok",
			StatusDelivery: "ok",
		},
	}

	userId := uint64(32)
	count := 3
	offset := 0
	sortString := ""
	paginator := models.PaginatorOrders{
		PageNum:           1,
		Count:             count,
		SortOrdersOptions: models.SortOrdersOptions{
			SortKey:       "",
			SortDirection: "",
		},
	}

	t.Run("SelectRangeOrders_internal_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "address", "total_cost", "date_added",
				"date_delivery", "order_num", "status_pay"}).
			AddRow(orders[0].Id, orders[0].Address.Address, orders[0].TotalCost,
			orders[0].DateAdded, orders[0].DateDelivery, orders[0].OrderNumber.Number,
			orders[0].StatusPay)
		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(userId, count, offset).
			WillReturnRows(rows)

		_, err = repository.SelectRangeOrders(userId, sortString, &paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("SelectRangeOrders_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "address", "total_cost", "date_added",
				"date_delivery", "order_num", "status_pay", "status_delivery"}).
			AddRow(orders[0].Id, orders[0].Address.Address, orders[0].TotalCost,
				orders[0].DateAdded, orders[0].DateDelivery, orders[0].OrderNumber.Number,
				orders[0].StatusPay, orders[0].StatusDelivery)
		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(userId, count, offset).
			WillReturnRows(rows)

		data, err := repository.SelectRangeOrders(userId, sortString, &paginator)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, orders, data, "data not equal")
	})
}

func TestPostgresqlRepository_GetCountPages(t *testing.T) {
	countId := 12
	userId := uint64(32)
	countOnPage := 1

	t.Run("GetCountPages_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"count(id)"}).
			AddRow(countId)
		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnRows(rows)

		data, err := repository.GetCountPages(userId, countOnPage)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, countId, data, "data not equal")
	})

	t.Run("GetCountPages_internal_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"count(id)"}).
			AddRow(nil)
		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(userId).
			WillReturnRows(rows)

		_, err = repository.GetCountPages(userId, countOnPage)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_CreateSortString(t *testing.T) {

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		data, err := repository.CreateSortString("date", "ASC")
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, "ORDER BY date_added ASC ", data, "data not equal")
	})

	t.Run("CreateSortString_success", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		data, err := repository.CreateSortString("date", "DESC")
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, "ORDER BY date_added DESC ", data, "data not equal")
	})
}
