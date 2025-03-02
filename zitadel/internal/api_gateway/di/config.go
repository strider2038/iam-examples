package di

import (
	"net/url"

	"github.com/caarlos0/env/v11"
	"github.com/muonsoft/errors"
)

type Config struct {
	ZitadelURL         url.URL `env:"ZITADEL_URL,required,notEmpty"`
	ZitadelClientID    string  `env:"ZITADEL_CLIENT_ID,required,notEmpty"`
	ZitadelRedirectURI string  `env:"ZITADEL_REDIRECT_URI,required,notEmpty"`
	EncryptionKey      string  `env:"ENCRYPTION_KEY,required,notEmpty"`
	TargetHost         string  `env:"TARGET_HOST,required,notEmpty"`
}

func ParseConfig() (Config, error) {
	config, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, errors.Errorf("parse env: %w", err)
	}

	return config, nil
}
