package usecase

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	Repository *repository.ProductRepository
	DB         *gorm.DB
}

func NewProductUseCase(Repository *repository.ProductRepository, DB *gorm.DB) *ProductUseCase {
	return &ProductUseCase{
		Repository: Repository,
		DB:         DB,
	}
}

func stringToProductStatus(statusStr string) model.ProductStatus {
	switch strings.ToLower(strings.TrimSpace(statusStr)) {
	case "active":
		return model.StatusActive
	case "inactive":
		return model.StatusInactive
	case "draft":
		return model.StatusDraft
	default:
		return model.StatusActive
	}
}

func (uc *ProductUseCase) ReadProductExcel(rows [][]string) []model.ProductImport {
	var products []model.ProductImport

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 5 {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			continue
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
		if err != nil {
			continue
		}

		product := model.ProductImport{
			Name:   name,
			Price:  price,
			Status: model.StatusActive,
		}

		if sku := strings.TrimSpace(row[0]); sku != "" {
			product.SKU = &sku
		}
		if len(row) > 2 {
			if slug := strings.TrimSpace(row[2]); slug != "" {
				product.Slug = &slug
			}
		}
		if len(row) > 3 {
			if desc := strings.TrimSpace(row[3]); desc != "" {
				product.Description = &desc
			}
		}
		if len(row) > 5 {
			if discountPrice, err := strconv.ParseFloat(strings.TrimSpace(row[5]), 64); err == nil {
				product.DiscountPrice = &discountPrice
			}
		}
		if len(row) > 6 {
			if stock, err := strconv.ParseInt(strings.TrimSpace(row[6]), 10, 64); err == nil {
				product.Stock = stock
			}
		}
		if len(row) > 7 {
			if weight, err := strconv.ParseFloat(strings.TrimSpace(row[7]), 64); err == nil {
				product.Weight = &weight
			}
		}
		if len(row) > 8 {
			product.Status = stringToProductStatus(row[8])
		}
		if len(row) > 9 {
			if isFeatured, err := strconv.ParseBool(strings.TrimSpace(row[9])); err == nil {
				product.IsFeatured = isFeatured
			}
		}

		products = append(products, product)
	}

	return products
}

func (uc *ProductUseCase) CreateProducts(c context.Context, products []model.ProductImport) error {
	tx := uc.DB.WithContext(c).Begin()

	defer tx.Rollback()

	var newProducts []model.Product

	for _, product := range products {
		if product.SKU != nil {
			err := uc.Repository.FindBySKU(c, tx, &model.Product{}, *product.SKU)

			// if product sku already exist, skip this product
			if err == nil {
				log.Printf("product sku already exist: %s", *product.SKU)
				continue
			}

			// other error besides not found from Repository.FindBySKU
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("unexpected error: %v", err)
				continue
			}
		}

		newProducts = append(newProducts, model.Product{
			SKU:           product.SKU,
			Name:          product.Name,
			Slug:          product.Slug,
			Description:   product.Description,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Stock:         product.Stock,
			Weight:        product.Weight,
			Status:        string(product.Status),
			IsFeatured:    product.IsFeatured,
		})
	}

	err := uc.Repository.Create(c, tx, &newProducts)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
