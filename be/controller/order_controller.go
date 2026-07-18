package controller

import (
	"fmt"
	"net/http"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/httputil"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/validator"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/worker"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
)

type OrderController struct {
	usecase      *usecase.OrderUseCase
	uploadWorker *worker.UploadWorker
	mqRoutingKey *configs.RabbitMQRoutingKey
	sftpPath     string
}

func NewOrderController(usecase *usecase.OrderUseCase, uploadWorker *worker.UploadWorker, mqRoutingKey *configs.RabbitMQRoutingKey, sftpPath string) *OrderController {
	return &OrderController{
		usecase:      usecase,
		uploadWorker: uploadWorker,
		mqRoutingKey: mqRoutingKey,
		sftpPath:     sftpPath,
	}
}

func (h *OrderController) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	paginationReq := httputil.GetPaginationParams(r)

	orderFilter := bindOrderFilterFromRequest(r)

	orders, err := h.usecase.FindAllOrders(r.Context(), paginationReq, orderFilter)

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", orders)

}

func (h *OrderController) UploadOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")

	if err != nil {
		model.WriteError(500, fmt.Sprintf("failed to read file %s", err.Error()))
		return
	}

	defer file.Close()

	isFileExtensionCorrect := validator.IsAllowedExtension(header.Filename)

	if !isFileExtensionCorrect {
		model.WriteError(400, "invalid file extension")
		return
	}

	h.uploadWorker.Queue(worker.UploadTask{
		File:            file,
		Filename:        header.Filename,
		Filepath:        h.sftpPath,
		QueueRoutingKey: h.mqRoutingKey.OrderDirectImport,
	})

	httputil.SendJsonResponse(w, 200, "success", nil)
}

func bindOrderFilterFromRequest(r *http.Request) *model.OrderFilter {
	return &model.OrderFilter{
		Email:       r.URL.Query().Get("email"),
		ProductName: r.URL.Query().Get("product_name"),
		Search:      r.URL.Query().Get("search"),
	}
}
