package repository

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresqlRepository_SelectRangeOrders(t *testing.T) {
	orders := []*models.PlacedOrder{
		{
			Id: uint64(1),
			Address: models.OrderAddress{
				Address: "test street",
			},
			TotalCost:    12,
			Products:     nil,
			DateAdded:    time.Time{},
			DateDelivery: time.Time{},
			OrderNumber:  models.OrderNumber{},
		},
	}

	userId := uint64(32)
	count := 3
	offset := 0
	sortString := ""
	paginator := models.PaginatorOrders{
		PageNum: 1,
		Count:   count,
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
				"date_delivery", "status"}).
			AddRow(orders[0].Id, orders[0].Address.Address, orders[0].TotalCost,
				orders[0].DateAdded, orders[0].DateDelivery,
				orders[0].OrderNumber.Number)
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
				"date_delivery", "order_num", "status"}).
			AddRow(orders[0].Id, orders[0].Address.Address, orders[0].TotalCost,
				orders[0].DateAdded, orders[0].DateDelivery,
				orders[0].OrderNumber.Number, orders[0].Status)
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

func TestPostgresqlRepository_GetProductsInOrder(t *testing.T) {
	previewOrders := models.PreviewOrderedProducts{}
	orderId := uint64(3)

	t.Run("GetProductsInOrder_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		rows := sqlmock.
			NewRows([]string{"id", "images"}).
			AddRow(previewOrders.Id, previewOrders.PreviewImage)
		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(orderId).
			WillReturnRows(rows)

		_, err = repository.GetProductsInOrder(orderId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetProductsInOrder_internal_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.
			ExpectQuery("SELECT").
			WithArgs(orderId).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.GetProductsInOrder(orderId)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_AddOrder(t *testing.T) {
	order := models.Order{}
	price := models.TotalPrice{}
	userId := uint64(3)
	orderId := uint64(2)
	orderNum := "0000-00000000"
	products := []*models.PreviewCartArticle{
		{},
	}

	t.Run("AddOrder_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		rows1 := sqlmock.
			NewRows([]string{"id", "order_num"}).
			AddRow(orderId, orderNum)
		rows2 := sqlmock.
			NewRows([]string{"id"}).
			AddRow(orderId)

		sqlMock.
			ExpectQuery("INSERT INTO user_orders").
			WithArgs(userId, order.Recipient.FirstName, order.Recipient.LastName,
				order.Recipient.Email, order.Address.Address, price.TotalBaseCost,
				price.TotalCost, price.TotalDiscount).
			WillReturnRows(rows1)

		sqlMock.
			ExpectQuery("INSERT INTO ordered_products").
			WithArgs(products[0].Id, orderId, products[0].Count,
				products[0].Price.BaseCost, products[0].Price.Discount).
			WillReturnRows(rows2)

		_, err = repository.AddOrder(&order, userId, products, &price)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddOrder_can't_insert_orders", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.
			ExpectQuery("INSERT INTO user_orders").
			WithArgs(userId, order.Recipient.FirstName, order.Recipient.LastName,
				order.Recipient.Email, order.Address.Address, price.TotalBaseCost,
				price.TotalCost, price.TotalDiscount).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.AddOrder(&order, userId, products, &price)
		assert.Error(t, err, "expected error")
	})

	t.Run("AddOrder_can't_insert_ordered_products", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		rows1 := sqlmock.
			NewRows([]string{"id", "order_num"}).
			AddRow(orderId, orderNum)

		sqlMock.
			ExpectQuery("INSERT INTO user_orders").
			WithArgs(userId, order.Recipient.FirstName, order.Recipient.LastName,
				order.Recipient.Email, order.Address.Address, price.TotalBaseCost,
				price.TotalCost, price.TotalDiscount).
			WillReturnRows(rows1)

		sqlMock.
			ExpectQuery("INSERT INTO ordered_products").
			WithArgs(products[0].Id, orderId, products[0].Count,
				products[0].Price.BaseCost, products[0].Price.Discount).
			WillReturnError(errors.ErrDBInternalError)

		_, err = repository.AddOrder(&order, userId, products, &price)
		assert.Error(t, err, "expected error")
	})
}

func TestPostgresqlRepository_ChangeStatusOrder(t *testing.T) {
	userId := uint64(3)
	orderId := uint64(2)
	orderNum := "0000-00000000"
	status := "доставлено"

	t.Run("ChangeStatusOrder_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)
		rows1 := sqlmock.
			NewRows([]string{"user_id", "order_num"}).
			AddRow(userId, orderNum)

		sqlMock.
			ExpectQuery("UPDATE user_orders").
			WithArgs(status, orderId).
			WillReturnRows(rows1)

		_, _, err = repository.ChangeStatusOrder(orderId, status)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("ChangeStatusOrder_internal_error", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		sqlMock.
			ExpectQuery("UPDATE user_orders").
			WithArgs(status, orderId).
			WillReturnError(errors.ErrDBInternalError)

		_, _, err = repository.ChangeStatusOrder(orderId, status)
		assert.Error(t, err, "expected error")
	})
}
