package repository

import (
	"context"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (r *UserRepository) CountAll(c context.Context) (*int64, error) {
	var totalRecords int64

	baseQuery := r.DB.WithContext(c).Model(&model.User{})

	err := baseQuery.Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	return &totalRecords, nil

}

func (r *UserRepository) FindAll(c context.Context, paginationReq *model.PaginationRequest, filter *model.UserFilter) (*model.Paginated[model.User], error) {
	var users []model.User
	var totalRecords int64

	offset := (paginationReq.Page - 1) * paginationReq.PerPage

	baseQuery := r.DB.WithContext(c).Model(&model.User{})

	query := filterUserQuery(filter, baseQuery)

	err := query.Session(&gorm.Session{}).Count(&totalRecords).Error

	if err != nil {
		return nil, err
	}

	err = query.Offset(offset).Limit(paginationReq.PerPage).Find(&users).Error

	totalPages := (int(totalRecords) + paginationReq.PerPage - 1) / paginationReq.PerPage

	if err != nil {
		return nil, err
	}

	return &model.Paginated[model.User]{
		Data:       users,
		Total:      totalRecords,
		TotalPages: totalPages,
	}, nil
}

func filterUserQuery(filter *model.UserFilter, query *gorm.DB) *gorm.DB {

	if filter.Search != "" {
		query = query.Where(
			"LOWER(name) ILIKE LOWER(?) OR LOWER(email) ILIKE LOWER(?) OR LOWER(phone_number) ILIKE LOWER(?)",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}
	if filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", filter.PhoneNumber)
	}

	return query
}
