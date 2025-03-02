package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ProxyHandler struct {
	targetHost string
	transport  http.RoundTripper
}

func NewProxyHandler(
	targetHost string,
	transport http.RoundTripper,
) http.Handler {
	targetHost = strings.TrimSuffix(targetHost, "/")

	return &ProxyHandler{
		targetHost: targetHost,
		transport:  transport,
	}
}

func (h *ProxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	targetURL := h.targetHost + request.URL.String()
	proxyRequest, err := http.NewRequestWithContext(request.Context(), request.Method, targetURL, request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("create proxy request: %+v", err), http.StatusInternalServerError)

		return
	}

	proxyRequest.ContentLength = request.ContentLength
	proxyRequest.Header = request.Header

	response, err := h.transport.RoundTrip(proxyRequest)
	if err != nil {
		http.Error(writer, fmt.Sprintf("send proxy request: %+v", err), http.StatusInternalServerError)

		return
	}
	defer response.Body.Close()

	for name, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(name, value)
		}
	}

	writer.WriteHeader(response.StatusCode)
	if _, err := io.Copy(writer, response.Body); err != nil {
		http.Error(writer, fmt.Sprintf("copy response body: %+v", err), http.StatusInternalServerError)

		return
	}
}
