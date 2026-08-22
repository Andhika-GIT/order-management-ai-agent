package controller

import (
	"net/http"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/httputil"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
)

type DashboardController struct {
	productUseCase *usecase.ProductUseCase
}

func NewDashboardController(productUseCase *usecase.ProductUseCase) *DashboardController {
	return &DashboardController{
		productUseCase: productUseCase,
	}
}

func (h *DashboardController) GetDataSummary(w http.ResponseWriter, r *http.Request) {
	totalProducts, err := h.productUseCase.CountAllProducts(r.Context())

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	response := &model.DashboardResponse{
		TotalProducts: *totalProducts,
	}

	httputil.SendJsonResponse(w, 200, "successfully get dashboard data", response)
}
