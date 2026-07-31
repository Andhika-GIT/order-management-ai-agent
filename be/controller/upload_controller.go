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
	s3Helper     S3_helper.S3Helper
	orderUsecase usecase.OrderUseCase
}

func NewUploadController(s3Helper S3_helper.S3Helper, orderUsecase usecase.OrderUseCase) *UploadController {
	return &UploadController{
		s3Helper:     s3Helper,
		orderUsecase: orderUsecase,
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

func (h *UploadController) GetUploadOrderImagePresignedURL(w http.ResponseWriter, r *http.Request) {
	var req model.OrderImagePresignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid request body"), nil)
		return
	}

	if !validator.IsAllowedImageExtension(req.Filename) {
		httputil.SendJsonErrorResponse(w, model.WriteError(400, "invalid image extension"), nil)
		return
	}

	order, err := h.orderUsecase.FindOrderByID(r.Context(), req.OrderID)

	if err != nil {
		errMsg := fmt.Sprintf("error when find order: %v", err)
		httputil.SendJsonErrorResponse(w, model.WriteError(400, errMsg), nil)
		return
	}

	key := fmt.Sprintf("uploads/%d/%d-%s", order.ID, time.Now().UnixNano(), req.Filename)

	url, err := h.s3Helper.GetPresignedInsertURL(r.Context(), key, req.ContentType, presignExpiry)

	if err != nil {
		httputil.SendJsonErrorResponse(w, model.WriteError(500, fmt.Sprintf("failed to generate presigned url: %s", err.Error())), nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", model.PresignResponse{URL: url, Key: key})
}
