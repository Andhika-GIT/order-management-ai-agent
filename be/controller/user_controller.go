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

type UserController struct {
	usecase      *usecase.UserUseCase
	uploadWorker *worker.UploadWorker
	mqRoutingKey *configs.RabbitMQRoutingKey
	sftpPath     string
}

func NewUserController(usecase *usecase.UserUseCase, uploadWorker *worker.UploadWorker, mqRoutingKey *configs.RabbitMQRoutingKey, sftpPath string) *UserController {
	return &UserController{
		usecase:      usecase,
		uploadWorker: uploadWorker,
		mqRoutingKey: mqRoutingKey,
		sftpPath:     sftpPath,
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
		QueueRoutingKey: h.mqRoutingKey.UserDirectImport,
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
