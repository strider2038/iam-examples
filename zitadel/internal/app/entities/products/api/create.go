package api

import (
	"fmt"
	"net/http"

	"app/internal/app/entities/products"
	"app/internal/pkg/auth"
)

type CreateHandler struct {
	useCase *products.CreateUseCase
}

func NewCreateHandler(useCase *products.CreateUseCase) *CreateHandler {
	return &CreateHandler{useCase: useCase}
}

func (h *CreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusBadRequest)

		return
	}

	session, err := auth.ParseSessionFromRequest(request)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusUnauthorized)

		return
	}

	command := products.CreateCommand{
		UserID:    session.UserID,
		CompanyID: session.CompanyID,
		Name:      request.FormValue("name"),
	}

	if _, err := h.useCase.Handle(request.Context(), command); err != nil {
		http.Error(writer, fmt.Sprintf("%+v", err), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
