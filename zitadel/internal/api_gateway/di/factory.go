package di

import (
	"context"
	"net/http"

	"app/internal/api_gateway/api"

	"github.com/muonsoft/errors"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func NewRouter(ctx context.Context, config Config) (*http.ServeMux, error) {
	proxyHandler := api.NewProxyHandler(config.TargetHost, http.DefaultTransport)

	options := make([]zitadel.Option, 0)
	if config.ZitadelURL.Scheme == "http" {
		options = append(options, zitadel.WithInsecure(config.ZitadelURL.Port()))
	}

	// Initiate the authentication by providing a zitadel configuration and handler.
	// This example will use OIDC/OAuth2 PKCE Flow, therefore you will also need to initialize that with the generated client_id:
	z := zitadel.New(config.ZitadelURL.Hostname(), options...)
	zitadelSDK, err := authentication.New(ctx,
		z, config.EncryptionKey,
		openid.DefaultAuthentication(
			config.ZitadelClientID,
			config.ZitadelRedirectURI,
			config.EncryptionKey,
			oidc.ScopeOpenID,
			// important! it is used to pass OrgID in the claims
			"urn:zitadel:iam:user:resourceowner",
		),
	)
	if err != nil {
		return nil, errors.Errorf("initialize zitadel sdk: %w", err)
	}

	// Initialize the middleware by providing the sdk
	mw := authentication.Middleware(zitadelSDK)

	router := http.NewServeMux()
	// Register the authentication handler on your desired path.
	// It will register the following handlers on it:
	// - /login (starts the authentication process to the Login UI)
	// - /callback (handles the redirect back from the Login UI)
	// - /logout (handles the logout process)
	router.Handle("/auth/", zitadelSDK)

	router.Handle("/", mw.RequireAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		// Using the [middleware.Context] function we can gather information about the authenticated user.
		authCtx := mw.Context(request.Context())
		request.Header.Set("X-User-Id", authCtx.UserInfo.Subject)

		if orgID := getOrgID(authCtx.UserInfo); orgID != "" {
			request.Header.Set("X-Company-Id", orgID)
		}

		proxyHandler.ServeHTTP(w, request)
	})))

	return router, nil
}

func getOrgID(info *oidc.UserInfo) string {
	if orgID, ok := info.Claims["urn:zitadel:iam:user:resourceowner:id"]; ok {
		if id, ok := orgID.(string); ok {
			return id
		}
	}

	return ""
}
