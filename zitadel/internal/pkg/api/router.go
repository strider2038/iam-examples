package api

import (
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func NewRouter(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range routes {
		pattern := strings.TrimSpace(route.Method + " " + route.Path)

		mux.Handle(pattern, route.Handler)
	}

	return mux
}
