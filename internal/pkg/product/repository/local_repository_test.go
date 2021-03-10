package repository

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

const goodProductId = 1
const goodPageNum = 1
const goodCount = 10

func TestLocalRepository_GetById(t *testing.T) {
	retProduct := &models.Product{
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
	}

	rep := NewSessionLocalRepository()
	res, err := rep.GetById(goodProductId)
	require.NoError(t, err)
	require.Equal(t, res, retProduct)
}
