package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
	"gorm.io/gorm"
)

type UserUseCase struct {
	Repository *repository.UserRepository
	DB         *gorm.DB
}

func NewUserUseCase(Repository *repository.UserRepository, DB *gorm.DB) *UserUseCase {
	return &UserUseCase{
		Repository: Repository,
		DB:         DB,
	}
}

func (uc *UserUseCase) ReadUsersExcel(rows [][]string) []model.UserImport {
	var users []model.UserImport

	for i, row := range rows {

		if i == 0 {
			continue
		}

		if len(row) >= 3 {
			users = append(users, model.UserImport{
				Name:        row[0],
				Email:       row[1],
				PhoneNumber: row[2],
			})
		}
	}

	return users
}

func (uc *UserUseCase) CreateNewUsers(c context.Context, users []model.UserImport) error {
	tx := uc.DB.WithContext(c).Begin()

	defer tx.Rollback()

	var newUsers []model.User

	for _, user := range users {

		err := uc.Repository.FindByEmail(c, tx, &model.User{}, user.Email)

		// if user email already exist, skip this user
		if err == nil {
			log.Printf("user already exist")
			continue
		}

		// other error besides not found from Repository.FindByEmail
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("unexpected error: %v", err)
			continue
		}

		newUsers = append(newUsers, model.User{
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		})

	}

	err := uc.Repository.Create(c, tx, &newUsers)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (uc *UserUseCase) FindUserByEmail(c context.Context, userEmail string) (model.UserResponse, error) {
	var user model.User

	tx := uc.DB.WithContext(c)

	err := uc.Repository.FindByEmail(c, tx, &user, userEmail)

	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
