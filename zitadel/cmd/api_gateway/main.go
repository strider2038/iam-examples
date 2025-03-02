package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"app/internal/api_gateway/di"

	"github.com/muonsoft/errors"
)

func main() {
	config, err := di.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	router, err := di.NewRouter(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("server listening, press ctrl+c to stop")
	if err := http.ListenAndServe(":8089", router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
