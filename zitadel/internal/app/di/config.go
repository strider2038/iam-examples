package di

import (
	"net/url"

	"github.com/caarlos0/env/v11"
	"github.com/muonsoft/errors"
)

type Config struct {
	ZitadelURL     url.URL `env:"ZITADEL_URL,required,notEmpty"`
	ZitadelKeyPath string  `env:"ZITADEL_KEY_PATH,required,notEmpty"`
	DatabaseURL    string  `env:"DATABASE_URL"`
}

func ParseConfig() (Config, error) {
	config, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, errors.Errorf("parse env: %w", err)
	}

	if config.DatabaseURL == "" {
		config.DatabaseURL = "./var/products.db"
	}

	return config, nil
}
