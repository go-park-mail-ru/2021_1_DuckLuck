package models

type CategoriesCatalog struct {
	Id        uint64               `json:"id"`
	LevelName string               `json:"level_name"`
	Next      []*CategoriesCatalog `json:"next"`
}

var (
	CategoryCatalogLevels = 3
)
