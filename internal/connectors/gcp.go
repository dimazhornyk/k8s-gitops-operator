package connectors

import (
	"context"
	"diploma/internal/common"
	"fmt"

	"google.golang.org/api/iam/v1"
)

type gcp struct {
	projectID string
	accountID string
	service   *iam.Service
}

func NewGCP(projectID, accountID string) (GCP, error) {
	service, err := iam.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	return &gcp{
		projectID: projectID,
		accountID: accountID,
		service:   service,
	}, nil
}

func (g *gcp) CreateServiceAccount(config common.ServiceConfig) (*iam.ServiceAccount, error) {
	if _, err := g.service.Projects.ServiceAccounts.Delete(config.ServiceName).Do(); err != nil {
		return nil, err
	}

	request := &iam.CreateServiceAccountRequest{
		AccountId: g.accountID,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.ServiceName,
		},
	}
	account, err := g.service.Projects.ServiceAccounts.Create("projects/"+g.projectID, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Create: %v", err)
	}

	bindings := make([]*iam.Binding, 0, len(config.Permissions.GCP))
	for _, permission := range config.Permissions.GCP {
		bindings = append(bindings, &iam.Binding{
			Members: []string{"serviceAccount:" + account.Email},
			Role:    permission.Role,
		})
	}

	policyRequest := &iam.SetIamPolicyRequest{
		Policy: &iam.Policy{
			Bindings: bindings,
		},
	}

	_, err = g.service.Projects.ServiceAccounts.SetIamPolicy("serviceAccount:"+account.Email, policyRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.SetIamPolicy: %v", err)
	}

	return account, nil
}
