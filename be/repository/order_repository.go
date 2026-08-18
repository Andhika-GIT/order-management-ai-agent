package repository

import (
	"context"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: DB,
	}
}

func (r *OrderRepository) CountAll(c context.Context) (*int64, error) {
	var totalRecords int64

	baseQuery := r.DB.WithContext(c).Model(&model.Order{})

	err := baseQuery.Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	return &totalRecords, nil

}

func (r *OrderRepository) FindAll(c context.Context, paginationReq *model.PaginationRequest, filter *model.OrderFilter) (*model.Paginated[model.Order], error) {
	var orders []model.Order
	var totalRecords int64

	offset := (paginationReq.Page - 1) * paginationReq.PerPage

	baseQuery := r.DB.WithContext(c).Model(&model.Order{}).Preload("User")

	query := filterOrderQuery(filter, baseQuery)

	err := query.Session(&gorm.Session{}).Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	err = baseQuery.Offset(offset).Limit(paginationReq.PerPage).Find(&orders).Error

	if err != nil {
		return nil, err
	}

	totalPages := (int(totalRecords) + paginationReq.PerPage - 1) / paginationReq.PerPage

	return &model.Paginated[model.Order]{
		Data:       orders,
		Total:      totalRecords,
		TotalPages: totalPages,
	}, nil
}

func (r *OrderRepository) FindById(c context.Context, orderID int64) (*model.Order, error) {
	var order model.Order

	err := r.DB.WithContext(c).Preload("User").First(&order, orderID).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrderAttachmentKey(c context.Context, orderID int64, attachment_key string) error {
	return r.DB.WithContext(c).Model(&model.Order{}).Where("id = ?", orderID).Update("attachment_key", attachment_key).Error
}

func filterOrderQuery(filter *model.OrderFilter, query *gorm.DB) *gorm.DB {

	if filter.Search != "" {
		query = query.Joins("LEFT JOIN users ON users.id = orders.user_id").Where(
			"LOWER(email) ILIKE LOWER(?) OR LOWER(product_name) ILIKE LOWER(?)",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}

	if filter.Email != "" {
		query = query.Where("users.email = ?", filter.Email)
	}
	if filter.ProductName != "" {
		query = query.Where("orders.phone_number = ?", filter.ProductName)
	}

	return query
}
