package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/internal/app/entities/products"
	"app/internal/pkg/auth"
)

type FindHandler struct {
	useCase *products.FindUseCase
}

var _ http.Handler = (*FindHandler)(nil)

func NewFindHandler(useCase *products.FindUseCase) *FindHandler {
	return &FindHandler{useCase: useCase}
}

func (h *FindHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusBadRequest)

		return
	}

	session, err := auth.ParseSessionFromRequest(request)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusUnauthorized)

		return
	}

	query := products.FindQuery{
		CompanyID: session.CompanyID,
		Search:    request.FormValue("search"),
	}

	items, err := h.useCase.Handle(request.Context(), query)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	response := strings.Builder{}
	for _, item := range items {
		userName := ""
		if item.CreatedBy != nil {
			userName = item.CreatedBy.FirstName + " " + item.CreatedBy.LastName
		}
		response.WriteString(
			"<tr>" +
				"<td>" + item.Name + "</td>" +
				"<td>" + item.CreatedAt.Format(time.DateTime) + "</td>" +
				"<td>" + userName + "</td>" +
				"</tr>",
		)
	}

	_, _ = writer.Write([]byte(response.String()))
}
