package model

import "time"

type Product struct {
	ID            int64     `json:"id" gorm:"primary_key;column:id"`
	SKU           *string   `json:"sku" gorm:"column:sku"`
	Name          string    `json:"name" gorm:"column:name"`
	Slug          *string   `json:"slug" gorm:"column:slug"`
	Description   *string   `json:"description" gorm:"column:description"`
	Price         float64   `json:"price" gorm:"column:price"`
	DiscountPrice *float64  `json:"discount_price" gorm:"column:discount_price"`
	Stock         int64     `json:"stock" gorm:"column:stock"`
	Weight        *float64  `json:"weight" gorm:"column:weight"`
	Status        string    `json:"status" gorm:"column:status"`
	IsFeatured    bool      `json:"is_featured" gorm:"column:is_featured"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (a *Product) TableName() string {
	return "products"
}

type ProductStatus string

const (
	StatusActive   ProductStatus = "active"
	StatusInactive ProductStatus = "inactive"
	StatusDraft    ProductStatus = "draft"
)

type ProductImport struct {
	SKU           *string       `json:"sku"`
	Name          string        `json:"name"`
	Slug          *string       `json:"slug"`
	Description   *string       `json:"description"`
	Price         float64       `json:"price"`
	DiscountPrice *float64      `json:"discount_price"`
	Stock         int64         `json:"stock"`
	Weight        *float64      `json:"weight"`
	Status        ProductStatus `json:"status"`
	IsFeatured    bool          `json:"is_featured"`
}
