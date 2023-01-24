package connectors

import "diploma/internal/common"

type Kubernetes interface {
	CreateDeployment(conf common.ServiceConfig) error
	CreateService(conf common.ServiceConfig) error
}

type Github interface {
	GetFile(repo, path string) ([]byte, error)
}

type Istio interface {
	CreateGateway() error
	CreateVirtualService(config common.ServiceConfig) error
}

type Storage interface {
	SaveConfigHash(repo string, hash string) error
	GetConfigHash(repo string) string
}
