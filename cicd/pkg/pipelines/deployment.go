package pipelines

import (
	log "github.com/kataras/golog"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/hi-devops-io/hi/cicd/pkg/config"
	"encoding/json"
	"os"
)


type App struct{
	Name string	`json:"name"`
	Project string `json:"project"`
	Profile string `json:"profile"`
	DockerRegistry string `json:"docker_registry"`
	ImageTag string `json:"image_tag"`
	Port int32 `json:"port"`
	Ports []int32 `json:"ports"`
}

type Response struct {
	Message string
	Code int
	Data interface{}
}

func init() {

}


func int32Ptr(i int32) *int32 { return &i }

// @Title Deploy
// @Description deploy application
// @Param app
// @Return deployment result as string, error
func Deploy(pipeline *config.Pipeline) (string, error)  {
	deploymentsClient := KubeApi.clientSet.AppsV1beta1().Deployments(pipeline.Project)

	if "" == pipeline.ImageTag {
		pipeline.ImageTag = "latest"
	}
	if "" == pipeline.DockerRegistry {
		pipeline.DockerRegistry = "docker-registry.default.svc:5000"
	}
	if "" == pipeline.Profile {
		pipeline.Profile = "dev"
	}

	log.Debug(pipeline)

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: pipeline.App,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"pipeline": pipeline.App,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  pipeline.App,
							Image: pipeline.DockerRegistry + "/" + pipeline.Project + "/" + pipeline.App + ":" + pipeline.ImageTag,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name: "APP_PROFILES_ACTIVE",
									Value: pipeline.Profile,
								},
							},
						},
					},
				},
			},
		},
	}
	log.Debug(deployment)

	// Create Deployment
	log.Info("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}

	retVal := fmt.Sprintf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	log.Info(retVal)

	createService(pipeline)

	return retVal, err
}


func createService(pipeline *config.Pipeline) (*apiv1.Service, error)   {
	// create service

	serviceClient := KubeApi.clientSet.CoreV1().Services(pipeline.Project)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: pipeline.App,
			Labels: map[string]string{
				"pipeline": pipeline.App,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Protocol:   apiv1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.IntOrString{IntVal: 8080},
				},
			},
		},
	}

	var svc *apiv1.Service
	var err error
	svc, err = serviceClient.Create(service)
	log.Info(svc, err)

	return svc, err
}