package models

type CategoriesCatalog struct {
	Id   uint64               `json:"id" valid:"type(uint64)"`
	Name string               `json:"name" valid:"utfletter, stringlength(3|30)"`
	Next []*CategoriesCatalog `json:"next,omitempty" valid:"notnull"`
}
