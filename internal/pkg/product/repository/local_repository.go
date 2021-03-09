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
				Title: "Chupa Chups assorted caramel \"Mini\"",
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
					"/product/6024670804.jpg", "/product/6024670805.jpg"},
			},
			4: &models.Product{
				Id:    4,
				Title: "Electric Toothbrush Oral-B PRO 6000",
				Price: models.ProductPrice{
					BaseCost: 50,
					Discount: 5,
				},
				Rating: 2,
				Description: "Oral-B is the # 1 brand of toothbrushes recommended " +
					"by most dentists in the world! * Discover the Oral-B PRO 6000. " +
					"Smart Series Triumph! The Oral-B PRO 6000 Smart Series Triumph Toothbrush " +
					"features Bluetooth 4.0 to sync with the free Oral-B App. " +
					"Take your brushing to the next level in 2 minutes as " +
					"recommended by your dentist for superior cleansing and gum health.",
				Category: "Appliances",
				Images: []string{"/product/6023124975.jpg", "/product/6023125065.jpg",
					"/product/6023125066.jpg"},
			},
			5: &models.Product{
				Id:    5,
				Title: "Bosch VitaPower Serie 4 jug blender",
				Price: models.ProductPrice{
					BaseCost: 30,
					Discount: 0,
				},
				Rating: 5,
				Description: "Pro-Performance System: Optimal texture for smoothies, " +
					"even with frozen fruits and hard ingredients. Reliable Bosch motor: " +
					"1200 W power with engine speeds up to 30,000 rpm. Easy assembly and " +
					"cleaning: knife, bowl, lid are dishwasher safe, easy to install bowl " +
					"and store cable. High-quality assembly in our own Bosch factory: " +
					"durability, reliable construction and high-quality materials. " +
					"ProEdge Blades: Brilliant mixing results thanks to efficient " +
					"and durable stainless steel blades. Made in Germany.",
				Category: "Appliances",
				Images:   []string{"/product/6026466446.jpg", "/product/6043224204.jpg", "/product/6043224631.jpg"},
			},
			6: &models.Product{
				Id:    6,
				Title: "Hairdryer Remington Shine Therapy",
				Price: models.ProductPrice{
					BaseCost: 75,
					Discount: 0,
				},
				Rating: 5,
				Description: "The Shine Therapy Pro professional hairdryer guarantees " +
					"maximum shine and salon-quality styling thanks to triple technology " +
					"for increasing hair shine: Super Ionic technology, " +
					"advanced micro-conditioners and uniform, constant heat.",
				Category: "Appliances",
				Images:   []string{"/product/6014172356.jpg", "/product/6014172357.jpg", "/product/6014172363.jpg"},
			},
			7: &models.Product{
				Id:    7,
				Title: "Electric grill Tefal OptiGrill",
				Price: models.ProductPrice{
					BaseCost: 120,
					Discount: 10,
				},
				Rating: 5,
				Description: "Just put the meat on the grill, select the desired " +
					"mode and start the process - a special patented sensor detects " +
					"the thickness and quantity of steaks and optimizes the cooking time.",
				Category: "Appliances",
				Images: []string{"/product/6024159698.jpg", "/product/6024159699.jpg",
					"/product/6024159700.jpg"},
			},
			8: &models.Product{
				Id:    8,
				Title: "Household vacuum cleaner Samsung",
				Price: models.ProductPrice{
					BaseCost: 30,
					Discount: 0,
				},
				Rating:      4,
				Description: "1600/350 W, Blue, 1.3 L, On / Off, Microfilter",
				Category:    "Appliances",
				Images: []string{"/product/1023571463.jpg", "/product/1023571464.jpg",
					"/product/6022334384.jpg"},
			},
			9: &models.Product{
				Id:    9,
				Title: "Fairy Platinum Plus Dishwasher Capsules",
				Price: models.ProductPrice{
					BaseCost: 15,
					Discount: 0,
				},
				Rating: 4,
				Description: "For incredible cleanliness and brilliance of your cookware, " +
					"try Fairy Platinum Plus All-in-One. Fairy dishwasher capsules handle " +
					"even tough dirt. They also help remove plaque build-up over time, " +
					"restoring crockery to its original shine. Fairy Platinum Plus " +
					"Dishwasher Capsules with 3 Liquid Compartments dissolve quickly " +
					"even at low temperatures. Capsules suitable for short dishwashing " +
					"cycles protect glass and silver and prevent limescale build-up " +
					"thanks to the efficient rinse function.",
				Category: "Food",
				Images: []string{"/product/6024507661.jpg", "/product/6024507662.jpg",
					"/product/6024507659.jpg"},
			},
			10: &models.Product{
				Id:    10,
				Title: "Twix Minis chocolates",
				Price: models.ProductPrice{
					BaseCost: 10.60,
					Discount: 5,
				},
				Rating: 5,
				Description: "Twix minis chocolate bars are crispy biscuits, thick caramel " +
					"and great milk chocolate. Drinking tea with Twix with colleagues, friends, " +
					"or family is a great way to spend your free time.",
				Category: "Food",
				Images: []string{"/product/1015450862.jpg", "/product/1023564323.jpg",
					"/product/1028485113.jpg"},
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
