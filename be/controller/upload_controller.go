package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/S3_helper"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/httputil"
	"github.com/Andhika-GIT/go-message-broker-monorepo/pkg/validator"
	"github.com/Andhika-GIT/go-message-broker-monorepo/usecase"
)

const presignExpiry = 15 * time.Minute

type UploadController struct {
	s3Helper       S3_helper.S3Helper
	productUsecase usecase.ProductUseCase
}

func NewUploadController(s3Helper S3_helper.S3Helper, productUsecase usecase.ProductUseCase) *UploadController {
	return &UploadController{
		s3Helper:       s3Helper,
		productUsecase: productUsecase,
	}
}

func (h *UploadController) GetPresignedUploadURL(w http.ResponseWriter, r *http.Request) {
	var req model.PresignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid request body"), nil)
		return
	}

	if !validator.IsAllowedExtension(req.Filename) {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid file extension"), nil)
		return
	}

	key := fmt.Sprintf("uploads/%d-%s", time.Now().UnixNano(), req.Filename)

	url, err := h.s3Helper.GetPresignedInsertURL(r.Context(), key, req.ContentType, presignExpiry)

	if err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(500, fmt.Sprintf("failed to generate presigned url: %s", err.Error())), nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", model.PresignResponse{URL: url, Key: key})
}

func (h *UploadController) GetUploadProductImagePresignedURL(w http.ResponseWriter, r *http.Request) {
	var req model.ProductImagePresignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid request body"), nil)
		return
	}

	if !validator.IsAllowedImageExtension(req.Filename) {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid image extension"), nil)
		return
	}

	product, err := h.productUsecase.FindProductByID(r.Context(), req.ProductID)

	if err != nil {
		errMsg := fmt.Sprintf("error when find product: %v", err)
		httputil.SendJsonErrorResponse(w, model.WriteError(400, errMsg), nil)
		return
	}

	key := fmt.Sprintf("uploads/%d/%d-%s", product.ID, time.Now().UnixNano(), req.Filename)

	url, err := h.s3Helper.GetPresignedInsertURL(r.Context(), key, req.ContentType, presignExpiry)

	if err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(500, fmt.Sprintf("failed to generate presigned url: %s", err.Error())), nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", model.PresignResponse{URL: url, Key: key})
}

func (h *UploadController) ConfirmUploadProductImage(w http.ResponseWriter, r *http.Request) {
	var req model.ProductImageConfirmRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid request body"), nil)
		return
	}

	product, err := h.productUsecase.FindProductByID(r.Context(), req.ProductID)

	if err != nil {
		errMsg := fmt.Sprintf("error when find product: %v", err)
		httputil.SendJsonErrorResponse(w, model.WriteError(400, errMsg), nil)
		return
	}

	if product.ImageURL != nil {
		h.s3Helper.Delete(r.Context(), *product.ImageURL)
	}

	err = h.productUsecase.InsertProductImageURL(r.Context(), product.ID, req.Key)

	if err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(500, fmt.Sprintf("failed to save new product image : %s", err.Error())), nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", nil)

}
