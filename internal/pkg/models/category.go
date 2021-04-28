package models

// Node of categories tree
type CategoriesCatalog struct {
	Id   uint64               `json:"id"`
	Name string               `json:"name" valid:"utfletter, stringlength(3|30)"`
	Next []*CategoriesCatalog `json:"next,omitempty" valid:"notnull"`
}
