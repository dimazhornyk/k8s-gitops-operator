package common

import (
	"k8s.io/api/core/v1"
)

func RoutesToContainerPorts(routes []Route) []v1.ContainerPort {
	var ports []v1.ContainerPort
	for _, route := range routes {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(route.Port),
		})
	}
	return ports
}
