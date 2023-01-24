package main

import (
	"diploma/internal/common"
	"diploma/internal/connectors"
	"diploma/internal/logic"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func main() {
	c, err := buildContainer()
	if err != nil {
		panic(err)
	}

	if err := c.Invoke(func(logger *logrus.Entry, service logic.Service) {
		logger.Info("Starting application")
		if err := service.Start(); err != nil {
			logger.Fatal(err)
		}
	}); err != nil {
		panic(err)
	}
}

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
	if err := c.Provide(connectors.NewKubernetes); err != nil {
		return nil, err
	}
	if err := c.Provide(connectors.NewIstio); err != nil {
		return nil, err
	}
	if err := c.Provide(connectors.NewStorage); err != nil {
		return nil, err
	}
	if err := c.Provide(logic.NewService); err != nil {
		return nil, err
	}

	return c, nil
}
