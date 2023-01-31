package logic

import (
	"crypto/sha256"
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
	k8s     connectors.Kubernetes
	github  connectors.Github
	istio   connectors.Istio
	storage connectors.Storage
	logger  *logrus.Entry
}

func NewService(k8s connectors.Kubernetes, istio connectors.Istio, github connectors.Github, storage connectors.Storage, logger *logrus.Entry) Service {
	return &service{
		k8s:     k8s,
		istio:   istio,
		github:  github,
		storage: storage,
		logger:  logger,
	}
}

func (s service) Start() error {
	//if err := s.istio.CreateGateway(); err != nil {
	//	s.logger.Info("skip creating gateway")
	//}

	http.HandleFunc("/", s.handleGithubWebhook)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func (s service) handleConfigChange(repoFullname string) error {
	b, err := s.github.GetFile(repoFullname, "ops.yaml")
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}

	hash, err := s.getConfHash(b)
	if err != nil {
		return errors.Wrap(err, "failed to get config hash")
	}

	if s.storage.GetConfigHash(repoFullname) == hash {
		s.logger.Infof("config for %s is unchanged, skipping", repoFullname)
		return nil
	}

	if err := s.storage.SaveConfigHash(repoFullname, hash); err != nil {
		return errors.Wrap(err, "failed to update config hash")
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

func (s service) getConfHash(conf []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(conf); err != nil {
		return "", errors.Wrap(err, "failed to write to hash")
	}

	return string(h.Sum(nil)), nil
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

	if m["ref"] != "refs/heads/main" {
		s.logger.Info("skip non-main branch")

		return
	}

	repoFullname := m["repository"].(map[string]interface{})["full_name"].(string)

	if err := s.handleConfigChange(repoFullname); err != nil {
		s.logger.WithError(err).Error("failed to handleConfigChange request")
	}
}
