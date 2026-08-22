package repository

import (
	"context"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: DB,
	}
}

func (r *ProductRepository) CountAll(c context.Context) (*int64, error) {
	var totalRecords int64

	baseQuery := r.DB.WithContext(c).Model(&model.Product{})

	err := baseQuery.Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	return &totalRecords, nil

}

func (r *ProductRepository) FindAll(c context.Context, paginationReq *model.PaginationRequest, filter *model.ProductFilter) (*model.Paginated[model.Product], error) {
	var products []model.Product
	var totalRecords int64

	offset := (paginationReq.Page - 1) * paginationReq.PerPage

	baseQuery := r.DB.WithContext(c).Model(&model.Product{})

	query := filterProductQuery(filter, baseQuery)

	err := query.Session(&gorm.Session{}).Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	err = query.Offset(offset).Limit(paginationReq.PerPage).Find(&products).Error

	if err != nil {
		return nil, err
	}

	totalPages := (int(totalRecords) + paginationReq.PerPage - 1) / paginationReq.PerPage

	return &model.Paginated[model.Product]{
		Data:       products,
		Total:      totalRecords,
		TotalPages: totalPages,
	}, nil
}

func (r *ProductRepository) FindById(c context.Context, productID int64) (*model.Product, error) {
	var product model.Product

	err := r.DB.WithContext(c).First(&product, productID).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) UpdateProductImageURL(c context.Context, productID int64, image_url string) error {
	return r.DB.WithContext(c).Model(&model.Product{}).Where("id = ?", productID).Update("image_url", image_url).Error
}

func filterProductQuery(filter *model.ProductFilter, query *gorm.DB) *gorm.DB {

	if filter.Search != "" {
		query = query.Where(
			"LOWER(name) ILIKE LOWER(?) OR LOWER(sku) ILIKE LOWER(?)",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}

	if filter.SKU != "" {
		query = query.Where("sku = ?", filter.SKU)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}
