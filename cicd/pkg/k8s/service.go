package k8s

import (
	"k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt"
	"github.com/hidevopsio/hi/cicd/pkg/config"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func CreateService(pipeline *config.Pipeline) error {
	// create service

	serviceSpec := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: pipeline.App,
			Labels: map[string]string{
				"app": pipeline.App,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8080,
				},
			},
			Selector: map[string]string{
				"app": pipeline.App,
			},
		},
	}

	services := Client.ClientSet.CoreV1().Services(pipeline.Project)
	svc, err := services.Get(pipeline.App, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = services.Update(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to update service: %s", err)
		}
		log.Info("service updated")
	case errors.IsNotFound(err):
		_, err = services.Create(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to create service")
		}
		log.Info("service created")
	default:
		return fmt.Errorf("upexected error: %s", err)
	}
	return nil
}
