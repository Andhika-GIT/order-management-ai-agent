package repository

import (
	"context"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"gorm.io/gorm"
)

type OrderRepository struct{}

func (r *OrderRepository) Create(c context.Context, tx *gorm.DB, order *[]model.Order) error {
	err := tx.Create(&order).Error

	if err != nil {
		return err
	}

	return nil
}
