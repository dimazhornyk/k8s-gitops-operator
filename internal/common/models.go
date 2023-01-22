package common

import v1 "k8s.io/api/core/v1"

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
	ExtraEnv    v1.EnvVar    `json:"extraEnv,omitempty"`
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

type Route struct {
	Name   string `json:"name"`
	Scope  string `json:"scope"`
	Prefix string `json:"prefix"`
	Port   uint32 `json:"port"`
}

type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

func RoutesToContainerPorts(routes []Route) []v1.ContainerPort {
	var ports []v1.ContainerPort
	for _, route := range routes {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(route.Port),
		})
	}
	return ports
}
