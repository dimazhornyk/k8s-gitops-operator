package connectors

import (
	"context"
	"diploma/internal/common"
	"github.com/pkg/errors"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type istio struct {
	clientset *versioned.Clientset
}

func NewIstio() (Istio, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	restConfig, err := clientcmd.BuildConfigFromFlags(config.Host, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create k8s rest client")
	}

	cs, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create istio client")
	}

	return &istio{
		clientset: cs,
	}, nil
}

func (i istio) CreateVirtualService(config common.ServiceConfig) error {
	vs := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.ServiceName,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{"*"},
			Gateways: []string{
				"istio-gateway",
			},
			Http: []*networkingv1alpha3.HTTPRoute{
				{Match: []*networkingv1alpha3.HTTPMatchRequest{
					{Headers: map[string]*networkingv1alpha3.StringMatch{
						"Host": {
							MatchType: &networkingv1alpha3.StringMatch_Exact{
								Exact: config.ServiceName,
							},
						},
					}},
				}},
			},
		},
	}

	opts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
	}

	_, err := i.clientset.NetworkingV1alpha3().VirtualServices("default").Create(context.Background(), vs, opts)
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
