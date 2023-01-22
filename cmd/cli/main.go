package cli

import (
	"diploma/internal/common"
	"diploma/internal/connectors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func buildContainer() (*dig.Container, error) {
	c := dig.New()
	if err := c.Provide(func(config *common.Config) (logrus.Level, error) {
		return logrus.ParseLevel(config.LogLevel)
	}); err != nil {
		return nil, err
	}
	if err := c.Provide(common.NewLogger); err != nil {
		return nil, err
	}
	if err := c.Provide(common.NewConfig); err != nil {
		return nil, err
	}
	if err := c.Provide(connectors.NewGithub); err != nil {
		return nil, err
	}

	return c, nil
}
