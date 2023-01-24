package connectors

import (
	"context"
	"diploma/internal/common"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type k8s struct {
	clientset *kubernetes.Clientset
}

func NewKubernetes() (Kubernetes, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k8s{
		clientset: cs,
	}, nil
}

func (k k8s) CreateDeployment(conf common.ServiceConfig) error {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   conf.ServiceName,
			Labels: map[string]string{"app": conf.ServiceName},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": conf.ServiceName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": conf.ServiceName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  conf.ServiceName,
							Image: fmt.Sprintf("%s:%s", conf.Image.Repository, conf.Image.Tag),
							Ports: common.RoutesToContainerPorts(conf.Routes),
							// TODO: remove for production environment
							ImagePullPolicy: corev1.PullNever,
						},
					},
				},
			},
		},
	}
	opts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
	}

	_, err := k.clientset.AppsV1().Deployments("default").Create(context.Background(), deployment, opts)

	return err
}

func (k k8s) CreateService(conf common.ServiceConfig) error {
	ports := make([]corev1.ServicePort, len(conf.Routes))
	for i, route := range conf.Routes {
		ports[i] = corev1.ServicePort{
			Port:       int32(route.Port),
			TargetPort: intstr.FromInt(int(route.Port)),
		}
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: conf.ServiceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": conf.ServiceName},
			Ports:    ports,
		},
	}

	opts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
	}

	_, err := k.clientset.CoreV1().Services("default").Create(context.Background(), service, opts)

	return err
}
