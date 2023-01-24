package logic

import (
	"diploma/internal/common"
	"diploma/internal/connectors"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type service struct {
	k8s    connectors.Kubernetes
	github connectors.Github
	istio  connectors.Istio
	logger *logrus.Entry
}

func NewService(istio connectors.Istio, github connectors.Github, logger *logrus.Entry) Service {
	return &service{
		istio:  istio,
		github: github,
		logger: logger,
	}
}

func (s service) Start() error {
	if err := s.istio.CreateGateway(); err != nil {
		return errors.Wrap(err, "failed to create gateway")
	}

	http.HandleFunc("/", s.handleGithubWebhook)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func (s service) handle(repoFullname string) error {
	b, err := s.github.GetFile(repoFullname, "ops.yaml")
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}

	var conf common.ServiceConfig
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return errors.Wrap(err, "failed to unmarshal file")
	}

	if err := s.k8s.CreateDeployment(conf); err != nil {
		return errors.Wrap(err, "failed to create deployment")
	}

	if err := s.k8s.CreateService(conf); err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	if err := s.istio.CreateVirtualService(conf); err != nil {
		return errors.Wrap(err, "failed to create virtual service")
	}

	return nil
}

func (s service) handleGithubWebhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.WithError(err).Error("failed to read body")

		return
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		s.logger.WithError(err).Error("failed to unmarshal request")

		return
	}

	repoFullname := m["repository"].(map[string]interface{})["full_name"].(string)

	if err := s.handle(repoFullname); err != nil {
		s.logger.WithError(err).Error("failed to handle request")
	}
}
