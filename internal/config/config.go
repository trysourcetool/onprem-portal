package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

var Config *config

const (
	EnvLocal   = "local"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

type config struct {
	Env string `env:"ENV"`
	Jwt struct {
		Key string `env:"JWT_KEY"`
	}
	Postgres struct {
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DB       string `env:"POSTGRES_DB"`
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
	}
	Resend struct {
		APIKey string `env:"RESEND_API_KEY"`
	}
}

func Init() {
	cfg := new(config)
	envOpts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, envOpts); err != nil {
		log.Fatal("[INIT] config: ", err)
	}

	Config = cfg
}
