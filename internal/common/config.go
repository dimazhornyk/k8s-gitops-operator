package common

import (
	"github.com/caarlos0/env"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	GithubToken string `env:"GITHUB_TOKEN" envDefault:"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"`
}

func NewConfig() (*Config, error) {
	c := new(Config)
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
