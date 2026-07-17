package controller

import (
	"net/http"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/httputil"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
)

type DashboardController struct {
	userUseCase  *usecase.UserUseCase
	orderUseCase *usecase.OrderUseCase
}

func NewDashboardController(userUseCase *usecase.UserUseCase, orderUseCase *usecase.OrderUseCase) *DashboardController {
	return &DashboardController{
		userUseCase:  userUseCase,
		orderUseCase: orderUseCase,
	}
}

func (h *DashboardController) GetDataSummary(w http.ResponseWriter, r *http.Request) {
	totalUsers, err := h.userUseCase.CountAllUsers(r.Context())

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	totalOrders, err := h.orderUseCase.CountAllOrders(r.Context())

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	response := &model.DashboardResponse{
		TotalOrders: *totalOrders,
		TotalUsers:  *totalUsers,
	}

	httputil.SendJsonResponse(w, 200, "successfully get dashboard data", response)
}
