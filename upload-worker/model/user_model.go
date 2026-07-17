package model

import "time"

type User struct {
	ID          int64     `json:"id" gorm:"primary_key;column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Email       string    `json:"email" gorm:"column:email"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime"`
}

func (a *User) TableName() string {
	return "users"
}

type UserImport struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
