package model

import "time"

type Order struct {
	ID          int64     `json:"id" gorm:"primary_key;column:id"`
	UserId      int64     `json:"user_id" gorm:"column:user_id"`
	ProductName string    `json:"product_name" gorm:"column:product_name"`
	Quantity    int64     `json:"quantity" gorm:"column:quantity"`
	Status      string    `json:"status" gorm:"column:status"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime"`

	// user relationship
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
}

func (a *Order) TableName() string {
	return "orders"
}

type OrderImport struct {
	Email       string `json:"email"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
}

type OrderResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
}

type OrderFilter struct {
	Email       string
	ProductName string
	Search      string
}
