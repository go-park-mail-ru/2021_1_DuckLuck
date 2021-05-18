package models

import (
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
)

var (
	OrdersDateAddedSort = "date"

	OrdersPaginatorASC  = "ASC"
	OrdersPaginatorDESC = "DESC"
)

// Set of product with count uniq sets of this size
type RangeOrders struct {
	ListPreviewOrders []*PlacedOrder `json:"list_placed_orders" valid:"notnull"`
	MaxCountPages     int            `json:"max_count_pages"`
}

// Paginator for showing page of orders
type PaginatorOrders struct {
	PageNum int `json:"page_num"`
	Count   int `json:"count"`
	SortOrdersOptions
}

type SortOrdersOptions struct {
	SortKey       string `json:"sort_key" valid:"in(date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

func (p *PaginatorOrders) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	p.SortKey = sanitizer.Sanitize(p.SortKey)
	p.SortDirection = sanitizer.Sanitize(p.SortDirection)
}

// All order information
// This models saved in database
type Order struct {
	Recipient OrderRecipient `json:"recipient" valid:"notnull"`
	Address   OrderAddress   `json:"address" valid:"notnull"`
}

func (o *Order) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	o.Recipient.FirstName = sanitizer.Sanitize(o.Recipient.FirstName)
	o.Recipient.LastName = sanitizer.Sanitize(o.Recipient.LastName)
	o.Recipient.Email = sanitizer.Sanitize(o.Recipient.Email)
	o.Address.Address = sanitizer.Sanitize(o.Address.Address)
}

// All order information
// This model preview info for user
type PreviewOrder struct {
	Products  []*PreviewOrderedProducts `json:"products" valid:"notnull"`
	Recipient OrderRecipient            `json:"recipient" valid:"notnull"`
	Price     TotalPrice                `json:"price" valid:"notnull"`
	Address   OrderAddress              `json:"address" valid:"notnull"`
}

// Order address for delivery
type OrderAddress struct {
	Address string `json:"address" valid:"utfletter, stringlength(1|30)"`
}

// Info about order recipient
type OrderRecipient struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(1|30)"`
	Email     string `json:"email" valid:"email"`
}

type PlacedOrder struct {
	Id           uint64                    `json:"id"`
	Address      OrderAddress              `json:"address" valid:"notnull"`
	TotalCost    int                       `json:"total_cost"`
	Products     []*PreviewOrderedProducts `json:"product_images" valid:"notnull"`
	DateAdded    time.Time                 `json:"date_added"`
	DateDelivery time.Time                 `json:"date_delivery"`
	OrderNumber  OrderNumber               `json:"order_number"`
	Status       string                    `json:"status"`
}

type PreviewOrderedProducts struct {
	Id           uint64 `json:"id"`
	PreviewImage string `json:"preview_image" valid:"minstringlength(1)"`
}

type OrderNumber struct {
	Number string `json:"number"`
}
