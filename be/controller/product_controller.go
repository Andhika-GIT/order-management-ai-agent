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

type ProductController struct {
	usecase      *usecase.ProductUseCase
	uploadWorker *worker.UploadWorker
	config       *configs.Config
}

func NewProductController(usecase *usecase.ProductUseCase, uploadWorker *worker.UploadWorker, config *configs.Config) *ProductController {
	return &ProductController{
		usecase:      usecase,
		uploadWorker: uploadWorker,
		config:       config,
	}
}

func (h *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	paginationReq := httputil.GetPaginationParams(r)

	productFilter := bindProductFilterFromRequest(r)

	products, err := h.usecase.FindAllProducts(r.Context(), paginationReq, productFilter)

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", products)

}

func (h *ProductController) UploadProduct(w http.ResponseWriter, r *http.Request) {
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

	key := fmt.Sprintf("upload-temp/product/%d-%s", time.Now().UnixNano(), header.Filename)

	h.uploadWorker.Queue(worker.UploadTask{
		File:            bytes.NewReader(data),
		Key:             key,
		Bucket:          h.config.S3Config.Bucket,
		ContentType:     "application/octet-stream",
		QueueRoutingKey: h.config.RabbitMQRoutingKey.ProductImport,
	})

	httputil.SendJsonResponse(w, 200, "success", nil)
}

func bindProductFilterFromRequest(r *http.Request) *model.ProductFilter {
	return &model.ProductFilter{
		SKU:    r.URL.Query().Get("sku"),
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
	}
}
