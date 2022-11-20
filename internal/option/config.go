package option

import (
	"fmt"

	"github.com/arxdsilva/bravo/internal/clients/exchange"
	"github.com/arxdsilva/bravo/internal/http"
	"github.com/arxdsilva/bravo/internal/logger"
	"github.com/arxdsilva/bravo/internal/storage/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

const prefix = "APP"

type Config struct {
	HTTP     http.Config
	Log      logger.Config
	DB       postgres.Config
	Exchange exchange.Config
}

func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process(prefix, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &config, nil
}

func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
