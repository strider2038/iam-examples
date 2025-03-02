package frontend

import (
	_ "embed"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

var _ http.Handler = (*Handler)(nil)

//go:embed templates/index.html
var index []byte

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = writer.Write(index)
}
