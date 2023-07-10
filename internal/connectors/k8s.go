package connectors

import (
	"context"
	"diploma/internal/common"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
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
	deployment := appsv1.Deployment(conf.ServiceName, "default").
		WithLabels(map[string]string{"app": conf.ServiceName}).
		WithSpec(
			appsv1.DeploymentSpec().
				WithSelector(
					applymetav1.LabelSelector().
						WithMatchLabels(map[string]string{"app": conf.ServiceName})).
				WithTemplate(
					applycorev1.PodTemplateSpec().
						WithLabels(map[string]string{"app": conf.ServiceName}).
						WithSpec(
							applycorev1.PodSpec().
								WithContainers(&applycorev1.ContainerApplyConfiguration{
									Name:  common.Ptr(conf.ServiceName),
									Image: common.Ptr(fmt.Sprintf("%s:%s", conf.Image.Repository, conf.Image.Tag)),
									Ports: common.RoutesToContainerPorts(conf.Routes),
									// TODO: remove for production environment
									ImagePullPolicy: common.Ptr(corev1.PullNever),
								}).WithServiceAccountName(conf.ServiceName),
						),
				),
		)

	opts := metav1.ApplyOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		FieldManager: "application/apply-patch",
	}

	_, err := k.clientset.AppsV1().Deployments("default").Apply(context.Background(), deployment, opts)

	return err
}

func (k k8s) CreateService(conf common.ServiceConfig) error {
	ports := make([]*applycorev1.ServicePortApplyConfiguration, len(conf.Routes))
	for i, route := range conf.Routes {
		if route.Scope != common.Internal {
			continue
		}

		ports[i] = &applycorev1.ServicePortApplyConfiguration{
			Port:       common.Ptr(int32(route.Port)),
			TargetPort: common.Ptr(intstr.FromInt(int(route.Port))),
		}
	}

	service := applycorev1.Service(conf.ServiceName, "default").WithSpec(
		applycorev1.ServiceSpec().
			WithSelector(map[string]string{"app": conf.ServiceName}).
			WithPorts(ports...),
	)

	opts := metav1.ApplyOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		FieldManager: "application/apply-patch",
	}

	_, err := k.clientset.CoreV1().Services("default").Apply(context.Background(), service, opts)

	return err
}

func (k k8s) CreateServiceAccount(conf common.ServiceConfig, serviceAccountEmail string) error {
	serviceAcc := applycorev1.ServiceAccount(conf.ServiceName, "default").
		WithAnnotations(map[string]string{
			"iam.gke.io/gcp-service-account": serviceAccountEmail,
		})

	opts := metav1.ApplyOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		FieldManager: "application/apply-patch",
	}

	_, err := k.clientset.CoreV1().ServiceAccounts("default").Apply(context.Background(), serviceAcc, opts)

	return err
}
