package di

import (
	"context"
	"database/sql"
	"net/http"

	"app/internal/app/entities/products"
	products_api "app/internal/app/entities/products/api"
	"app/internal/app/entities/products/repository"
	user_repository "app/internal/app/entities/users/repository"
	"app/internal/app/frontend"
	"app/internal/pkg/api"

	"github.com/muonsoft/errors"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func NewRouter(ctx context.Context, config Config) (*http.ServeMux, error) {
	db, err := sql.Open("sqlite", config.DatabaseURL)
	if err != nil {
		return nil, errors.Errorf("open database: %w", err)
	}
	if err := initDB(db); err != nil {
		return nil, errors.Errorf("init database: %w", err)
	}

	options := make([]zitadel.Option, 0)
	if config.ZitadelURL.Scheme == "http" {
		options = append(options, zitadel.WithInsecure(config.ZitadelURL.Port()))
	}

	zitadelClient, err := client.New(ctx,
		zitadel.New(config.ZitadelURL.Hostname(), options...),
		client.WithAuth(client.DefaultServiceUserAuthentication(config.ZitadelKeyPath, oidc.ScopeOpenID, client.ScopeZitadelAPI())),
	)
	if err != nil {
		return nil, errors.Errorf("initialize zitadel sdk: %w", err)
	}

	userRepository := user_repository.NewZitadelUserRepository(zitadelClient)
	productRepository := repository.NewSQLiteProductRepository(db)
	productViewRepository := repository.NewProductViewRepository(productRepository, userRepository)

	router := api.NewRouter([]api.Route{
		{
			Method:  "GET",
			Path:    "/",
			Handler: frontend.NewHandler(),
		},
		{
			Method: "POST",
			Path:   "/api/products/find",
			Handler: products_api.NewFindHandler(
				products.NewFindUseCase(productViewRepository),
			),
		},
		{
			Method: "POST",
			Path:   "/api/products/create",
			Handler: products_api.NewCreateHandler(
				products.NewCreateUseCase(productRepository),
			),
		},
	})

	return router, nil
}
