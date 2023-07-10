package connectors

import (
	"diploma/internal/common"
	"google.golang.org/api/iam/v1"
)

type Kubernetes interface {
	CreateDeployment(conf common.ServiceConfig) error
	CreateService(conf common.ServiceConfig) error
	CreateServiceAccount(conf common.ServiceConfig, serviceAccountEmail string) error
}

type Github interface {
	GetFile(repo, path string) ([]byte, error)
}

type Istio interface {
	CreateVirtualService(config common.ServiceConfig) error
}

type Storage interface {
	SaveConfigHash(repo string, hash string) error
	GetConfigHash(repo string) string
}

type GCP interface {
	CreateServiceAccount(config common.ServiceConfig) (*iam.ServiceAccount, error)
}
