package usecase

import (
	"context"
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/rabbitmq"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"gorm.io/gorm"
)

type UserUseCase struct {
	Repository *repository.UserRepository
	rmq        *rabbitmq.RabbitMqProducer
	DB         *gorm.DB
}

func NewUserUseCase(Repository *repository.UserRepository, rmq *rabbitmq.RabbitMqProducer, DB *gorm.DB) *UserUseCase {
	return &UserUseCase{
		Repository: Repository,
		rmq:        rmq,
		DB:         DB,
	}
}

func (u *UserUseCase) CountAllUsers(c context.Context) (*int64, error) {
	totalUsers, err := u.Repository.CountAll(c)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to count all users %s", err.Error()))
	}

	return totalUsers, nil
}

func (u *UserUseCase) FindAllUsers(c context.Context, paginationReq *model.PaginationRequest, filter *model.UserFilter) (*model.Paginated[model.UserResponse], error) {

	paginated, err := u.Repository.FindAll(c, paginationReq, filter)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to find all users %s", err.Error()))
	}

	formatedUsers := convertToUsersResponse(paginated.Data)

	// return new paginated response with different type (UserResponse)
	return &model.Paginated[model.UserResponse]{
		Data:       formatedUsers,
		Total:      paginated.Total,
		TotalPages: paginated.TotalPages,
	}, nil

}

func convertToUsersResponse(users []model.User) []model.UserResponse {
	var usersResp []model.UserResponse

	for _, user := range users {
		resp := model.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		}

		usersResp = append(usersResp, resp)
	}

	return usersResp
}
