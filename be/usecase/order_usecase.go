package usecase

import (
	"context"
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
)

type OrderUseCase struct {
	Repository *repository.OrderRepository
}

func NewOrderUseCase(Repository *repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		Repository: Repository,
	}
}

func (u *OrderUseCase) CountAllOrders(c context.Context) (*int64, error) {
	totalOrders, err := u.Repository.CountAll(c)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to count all orders %s", err.Error()))
	}

	return totalOrders, nil
}

func (u *OrderUseCase) FindAllOrders(c context.Context, paginationReq *model.PaginationRequest, filter *model.OrderFilter) (*model.Paginated[model.OrderResponse], error) {
	paginated, err := u.Repository.FindAll(c, paginationReq, filter)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to find all users %s", err.Error()))
	}

	formatedOrders := convertToOrdersResponse(paginated.Data)

	// return new paginated response with different type (OrderResponse)
	return &model.Paginated[model.OrderResponse]{
		Data:       formatedOrders,
		Total:      paginated.Total,
		TotalPages: paginated.TotalPages,
	}, nil

}

func (u *OrderUseCase) FindOrderByID(c context.Context, orderID int64) (*model.OrderResponse, error) {
	order, err := u.Repository.FindById(c, orderID)

	if err != nil {
		return nil, err
	}

	resp := convertToOrderResponse(*order)

	return &resp, nil

}

func convertToOrdersResponse(orders []model.Order) []model.OrderResponse {
	var ordersResp []model.OrderResponse

	for _, order := range orders {
		resp := model.OrderResponse{
			ID:          order.ID,
			Email:       order.User.Email,
			ProductName: order.ProductName,
			Quantity:    order.Quantity,
		}

		ordersResp = append(ordersResp, resp)
	}

	return ordersResp
}

func convertToOrderResponse(order model.Order) model.OrderResponse {

	return model.OrderResponse{
		ID:          order.ID,
		Email:       order.User.Email,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
	}
}
