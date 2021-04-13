package repository

import (
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresqlRepository_AddOrder(t *testing.T) {
	orderId := uint64(2)
	productId := uint64(3)
	userId := uint64(4)
	order := models.Order{
		Recipient: models.OrderRecipient{
			FirstName: "test_name",
			LastName:  "test_last_name",
			Email:     "em@test.ru",
		},
		Address: models.OrderAddress{
			Address: "test1",
		},
	}
	price := models.TotalPrice{
		TotalDiscount: 10,
		TotalCost:     54,
		TotalBaseCost: 60,
	}
	products := []*models.PreviewCartArticle{
		&models.PreviewCartArticle{
			Id:    3,
			Title: "test",
			Price: models.ProductPrice{
				Discount: 3,
				BaseCost: 32,
			},
			PreviewImage: "tfdf",
			Count:        3,
		},
	}

	t.Run("SelectProductById_success", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		orderRows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
		sqlMock.ExpectQuery("INSERT INTO user_orders").
			WithArgs(userId, order.Recipient.FirstName, order.Recipient.LastName,
				order.Recipient.Email, order.Address.Address, price.TotalBaseCost,
				price.TotalCost, price.TotalDiscount).
			WillReturnRows(orderRows)

		productRows := sqlmock.NewRows([]string{"id"}).AddRow(productId)
		sqlMock.ExpectQuery("INSERT INTO ordered_products").
			WithArgs(products[0].Id, orderId, products[0].Count, products[0].Price.BaseCost,
				products[0].Price.Discount).
			WillReturnRows(productRows)

		savedId, err := repository.AddOrder(&order, userId, products, &price)
		if err != nil {
			t.Errorf("internal err: %s", err)
			return
		}

		err = sqlMock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met in order: %s", err)
			return
		}

		if savedId != orderId {
			t.Errorf("user wasn't saved: %s", err)
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

		sqlMock.ExpectQuery("INSERT INTO user_orders").
			WillReturnError(sql.ErrConnDone)

		_, err = repository.AddOrder(&order, userId, products, &price)
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

	t.Run("SelectProductById_bad_order", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("can't create mock: %s", err)
		}
		defer db.Close()

		repository := NewSessionPostgresqlRepository(db)

		orderRows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
		sqlMock.ExpectQuery("INSERT INTO user_orders").
			WithArgs(userId, order.Recipient.FirstName, order.Recipient.LastName,
				order.Recipient.Email, order.Address.Address, price.TotalBaseCost,
				price.TotalCost, price.TotalDiscount).
			WillReturnRows(orderRows)

		sqlMock.ExpectQuery("INSERT INTO ordered_products").
			WillReturnError(sql.ErrConnDone)

		_, err = repository.AddOrder(&order, userId, products, &price)
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
