package connectors

import (
	"context"
	"diploma/internal/common"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/applyconfiguration/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type istio struct {
	clientset *versioned.Clientset
	logger    *logrus.Entry
}

func NewIstio() (Istio, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create istio client")
	}

	return &istio{
		clientset: cs,
	}, nil
}

func (i istio) CreateVirtualService(config common.ServiceConfig) error {
	var externalRoute *common.Route
	for _, route := range config.Routes {
		if route.Scope == common.External {
			externalRoute = &route
			break
		}
	}

	if externalRoute == nil {
		i.logger.Info("skipping virtual service creation because there are no external routes")

		return nil
	}

	vs := clientnetworking.VirtualService(config.ServiceName, "default").
		WithSpec(networkingv1alpha3.VirtualService{
			Hosts: []string{"*"},
			Gateways: []string{
				"istio-gateway",
			},
			Http: []*networkingv1alpha3.HTTPRoute{
				{
					Name: fmt.Sprintf("%s-route", config.ServiceName),
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						{Headers: map[string]*networkingv1alpha3.StringMatch{
							"Host": {
								MatchType: &networkingv1alpha3.StringMatch_Exact{
									Exact: config.ServiceName,
								},
							},
						}},
					},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Host: config.ServiceName,
								Port: &networkingv1alpha3.PortSelector{
									Number: externalRoute.Port,
								},
							},
						},
					},
				},
			},
		})

	opts := metav1.ApplyOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		FieldManager: "application/apply-patch",
	}

	_, err := i.clientset.NetworkingV1alpha3().VirtualServices("default").Apply(context.Background(), vs, opts)
	if err != nil {
		return errors.Wrap(err, "Failed to create istio virtual service")
	}

	return nil
}

func (i istio) CreateGateway() error {
	gw := &v1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "istio-gateway",
		},

		Spec: networkingv1alpha3.Gateway{
			Selector: map[string]string{
				"istio": "ingressgateway",
			},
			Servers: []*networkingv1alpha3.Server{
				{
					Port: &networkingv1alpha3.Port{
						Number:   80,
						Protocol: "HTTP",
						Name:     "http",
					},
					Hosts: []string{
						"*",
					},
				},
			},
		},
	}

	opts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Gateway",
			APIVersion: "networking.istio.io/v1alpha3",
		},
	}

	_, err := i.clientset.NetworkingV1alpha3().Gateways("default").Create(context.Background(), gw, opts)
	if err != nil {
		return errors.Wrap(err, "Failed to create istio gateway")
	}

	return nil
}

//func (i istio) DeleteVirtualService(config common.ServiceConfig) error {
//	opts := metav1.DeleteOptions{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "VirtualService",
//			APIVersion: "networking.istio.io/v1alpha3",
//		},
//	}
//
//	if err := i.clientset.NetworkingV1alpha3().VirtualServices("default").Delete(context.Background(), config.ServiceName, opts); err != nil {
//		i.logger.WithError(err).Info("error deleting a virtual service")
//	}
//
//	return nil
//}
