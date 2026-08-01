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

type UserController struct {
	usecase      *usecase.UserUseCase
	uploadWorker *worker.UploadWorker
	config       *configs.Config
}

func NewUserController(usecase *usecase.UserUseCase, uploadWorker *worker.UploadWorker, config *configs.Config) *UserController {
	return &UserController{
		usecase:      usecase,
		uploadWorker: uploadWorker,
		config:       config,
	}
}

func (h *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	paginationReq := httputil.GetPaginationParams(r)

	userFilter := bindUserFilterFromRequest(r)

	users, err := h.usecase.FindAllUsers(r.Context(), paginationReq, userFilter)

	if err != nil {
		httputil.SendJsonErrorResponse(w, err, nil)
		return
	}

	httputil.SendJsonResponse(w, 200, "success", users)
}

func (h *UserController) UploadUser(w http.ResponseWriter, r *http.Request) {
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

	key := fmt.Sprintf("upload-temp/user/%d-%s", time.Now().UnixNano(), header.Filename)

	h.uploadWorker.Queue(worker.UploadTask{
		File:            bytes.NewReader(data),
		Key:             key,
		Bucket:          h.config.S3Config.Bucket,
		ContentType:     "application/octet-stream",
		QueueRoutingKey: h.config.RabbitMQRoutingKey.UserDirectImport,
	})

	httputil.SendJsonResponse(w, 200, "success", nil)
}

func bindUserFilterFromRequest(r *http.Request) *model.UserFilter {
	return &model.UserFilter{
		Name:        r.URL.Query().Get("name"),
		Email:       r.URL.Query().Get("email"),
		PhoneNumber: r.URL.Query().Get("phone_number"),
		Search:      r.URL.Query().Get("search"),
	}
}
