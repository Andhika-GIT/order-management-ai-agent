package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

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
	config       *configs.Config
}

func NewOrderController(usecase *usecase.OrderUseCase, uploadWorker *worker.UploadWorker, config *configs.Config) *OrderController {
	return &OrderController{
		usecase:      usecase,
		uploadWorker: uploadWorker,
		config:       config,
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
		httputil.SendJsonErrorResponse(w, model.WriteError(500, fmt.Sprintf("failed to read file %s", err.Error())), nil)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(500, "failed to read file"), nil)
		return
	}

	isFileExtensionCorrect := validator.IsAllowedExtension(header.Filename)

	if !isFileExtensionCorrect {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid file extension"), nil)
		return
	}

	key := fmt.Sprintf("upload-temp/order/%d-%s", time.Now().UnixNano(), header.Filename)

	h.uploadWorker.Queue(worker.UploadTask{
		File:            bytes.NewReader(data),
		Key:             key,
		Bucket:          h.config.S3Config.Bucket,
		ContentType:     "application/octet-stream",
		QueueRoutingKey: h.config.RabbitMQRoutingKey.OrderImport,
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
