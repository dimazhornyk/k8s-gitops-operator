package common

import "k8s.io/api/core/v1"

type Scope string

const (
	Internal Scope = "internal"
	External Scope = "external"
)

type RepositoryEvent struct {
	ID      string       `json:"id"`
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

type EventPayload struct {
	Ref  string `json:"ref"`
	Head string `json:"head"`
}

type ServiceConfig struct {
	ServiceName string       `json:"name"`
	Image       Image        `json:"image"`
	ExtraEnv    []v1.EnvVar  `json:"extraEnv,omitempty"`
	Routes      []Route      `json:"routes,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	GCP []GCPPermission `json:"gcp"`
}

type GCPPermission struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func (s Scope) String() string {
	return string(s)
}

func (s Scope) IsValid() bool {
	return s == Internal || s == External
}

type Route struct {
	Name  string `json:"name"`
	Scope Scope  `json:"scope"`
	Port  uint32 `json:"port"`
}

type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}
