package repository

import (
	"context"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"gorm.io/gorm"
)

type ProductRepository struct{}

func (r *ProductRepository) Create(c context.Context, tx *gorm.DB, products *[]model.Product) error {
	err := tx.Create(&products).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) FindBySKU(c context.Context, tx *gorm.DB, product *model.Product, sku string) error {
	err := tx.Where("sku = ?", sku).First(&product).Error

	if err != nil {
		return err
	}

	return nil
}
