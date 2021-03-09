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
				Title: "Hair dryer brush Rowenta",
				Price: models.ProductPrice{
					BaseCost: 20,
					Discount: 20,
				},
				Rating: 4,
				Description: "The rotating Brush Activ 'airstyler provides " +
					"unsurpassed drying results. Power of 1000" +
					"W guarantees fast drying effortlessly, two" +
					"rotating brushes with a diameter of 50 or 40 mm provide" +
					"professional styling. Ion generator and" +
					"ceramic coating smoothes hair, leaving it soft" +
					"and more brilliant.",
				Category: "Appliances",
				Images: []string{"/product/1021166584.jpg", "/product/1021166585.jpg",
					"/product/1021166586.jpg", "/product/6043447767.jpg"},
			},
			2: &models.Product{
				Id:    2,
				Title: "Chupa Chups assorted caramel \"Mini\", 100 pcs, 6 g each",
				Price: models.ProductPrice{
					BaseCost: 6.25,
					Discount: 0,
				},
				Rating: 3,
				Description: "Chupa Chups Mini is Chupa Chups' favorite candy on a stick " +
					"in mini format. In the showbox there are 100 Chupa. " +
					"Chups with the best flavors: strawberry, cola, orange, apple.",
				Category: "Food",
				Images: []string{"/product/6024670802.jpg", "/product/6024670803.jpg",
					"/product/6024670804.jpg", "/product/6024670805.jpg", "/product/6024670806.jpg"},
			},
			3: &models.Product{
				Id:    3,
				Title: "Iris Meller Chocolate, 24 pcs 38 g",
				Price: models.ProductPrice{
					BaseCost: 10.20,
					Discount: 10,
				},
				Rating: 5,
				Description: "Meller Chocolate is a true legend among sweets! " +
					"The combination of delicate caramel and chocolate, created using a " +
					"unique recipe, allows Meller to remain the leader in consumer preferences " +
					"for many years. A timeless classic and a truly legendary product!",
				Category: "Food",
				Images:   []string{"/product/6033457624.jpg", "/product/6033457625.jpg"},
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
			Price:        item.Price,
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
				return products[i].Price.BaseCost < products[j].Price.BaseCost
			}
		case models.PaginatorDESC:
			compare = func(i, j int) bool {
				return products[i].Price.BaseCost > products[j].Price.BaseCost
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
