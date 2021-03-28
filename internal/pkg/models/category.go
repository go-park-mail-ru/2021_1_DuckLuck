package models

type CategoriesCatalog struct {
	Id   uint64               `json:"id"`
	Name string               `json:"name"`
	Next []*CategoriesCatalog `json:"next,omitempty"`
}
