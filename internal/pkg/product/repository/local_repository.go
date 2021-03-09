package repository

import (
	"sort"
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[uint64]*models.Product
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() product.Repository {
	return &LocalRepository{
		data: map[uint64]*models.Product{
			1: &models.Product{
				Id:    1,
				Title: "Test1",
				Cost: models.ProductCost{
					BaseCost: 100,
					Discount: 20,
				},
				Rating:      1,
				Description: "Good item",
				Category:    "Home",
				Images:      []string{"/product/test.png"},
			},
			2: &models.Product{
				Id:    1,
				Title: "Test2",
				Cost: models.ProductCost{
					BaseCost: 50,
					Discount: 0,
				},
				Rating:      1,
				Description: "Good item",
				Category:    "Home",
				Images:      []string{"/product/test.png"},
			},
			3: &models.Product{
				Id:    1,
				Title: "Test1",
				Cost: models.ProductCost{
					BaseCost: 100,
					Discount: 20,
				},
				Rating:      1,
				Description: "Good item",
				Category:    "Home",
				Images:      []string{"/product/test.png"},
			},
		},
		mu: &sync.RWMutex{},
	}
}

func (lr *LocalRepository) GetById(productId uint64) (*models.Product, error) {
	lr.mu.RLock()
	productById, ok := lr.data[productId]
	lr.mu.RUnlock()

	if !ok {
		return nil, server_errors.ErrProductNotFound
	}

	return productById, nil
}

func (lr *LocalRepository) GetListPreviewProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, server_errors.ErrIncorrectPaginator
	}

	countPages := len(lr.data) / paginator.Count
	if len(lr.data)%paginator.Count > 0 {
		countPages++
	}

	if countPages < paginator.PageNum {
		return nil, server_errors.ErrProductsIsEmpty
	}

	products := make([]*models.ViewProduct, 0)
	lr.mu.RLock()
	for _, item := range lr.data {
		products = append(products, &models.ViewProduct{
			Id:           item.Id,
			Title:        item.Title,
			Cost:         item.Cost,
			Rating:       item.Rating,
			PreviewImage: item.Images[0],
		})
	}
	lr.mu.RUnlock()

	compare := func(i, j int) bool { return false }
	switch paginator.SortKey {
	case models.ProductsCostSort:
		switch paginator.SortDirection {
		case models.PaginatorASC:
			compare = func(i, j int) bool {
				return products[i].Cost.BaseCost < products[j].Cost.BaseCost
			}
		case models.PaginatorDESC:
			compare = func(i, j int) bool {
				return products[i].Cost.BaseCost > products[j].Cost.BaseCost
			}
		}
	case models.ProductsRatingSort:
		switch paginator.SortDirection {
		case models.PaginatorASC:
			compare = func(i, j int) bool {
				return products[i].Rating < products[j].Rating
			}
		case models.PaginatorDESC:
			compare = func(i, j int) bool {
				return products[i].Rating > products[j].Rating
			}
		}
	}
	sort.SliceStable(products, compare)

	leftBorder := (paginator.PageNum - 1) * paginator.Count
	rightBorder := paginator.PageNum * paginator.Count
	if paginator.PageNum == countPages {
		rightBorder = len(lr.data)
	}

	return &models.RangeProducts{
		ListPreviewProducts: products[leftBorder:rightBorder],
		MaxCountPages:       countPages,
	}, nil
}
