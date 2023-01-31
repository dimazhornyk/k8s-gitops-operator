package common

import (
	"k8s.io/client-go/applyconfigurations/core/v1"
)

func RoutesToContainerPorts(routes []Route) []v1.ContainerPortApplyConfiguration {
	var ports []v1.ContainerPortApplyConfiguration
	for _, route := range routes {
		ports = append(ports, v1.ContainerPortApplyConfiguration{
			ContainerPort: Ptr(int32(route.Port)),
		})
	}
	return ports
}

func Ptr[T any](val T) *T {
	return &val
}
