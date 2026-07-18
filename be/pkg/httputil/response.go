package httputil

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
)

func SendJsonResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := model.APIResponse{
		Status:  statusCode,
		Success: true,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(res)
}

func SendJsonErrorResponse(w http.ResponseWriter, err error, data any) {
	if e, ok := err.(*model.Error); ok {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(e.Code)

		res := model.APIResponse{
			Status:  e.Code,
			Success: false,
			Message: e.Message,
			Data:    data,
		}

		_ = json.NewEncoder(w).Encode(res)
	}
}

func GetPaginationParams(r *http.Request) *model.PaginationRequest {
	query := r.URL.Query()

	pageStr := query.Get("page")
	perPageStr := query.Get("per_page")

	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)

	if err != nil || page < 1 {
		perPage = 1
	}

	return &model.PaginationRequest{
		Page:    page,
		PerPage: perPage,
	}

}
