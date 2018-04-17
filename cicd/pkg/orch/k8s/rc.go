package k8s

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"fmt"
)

type ReplicationController struct{
	Name string
	Namespace string

	Interface v1.ReplicationControllerInterface
}

func NewReplicationController(name string, namespace string) *ReplicationController {
	return &ReplicationController{
		Name: name,
		Namespace: namespace,
		Interface: ClientSet.CoreV1().ReplicationControllers(namespace),
	}
}


func (rc *ReplicationController) Create(replicas int32) (*corev1.ReplicationController, error) {
	crc := &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Name: rc.Name,
			Labels: map[string]string{
				"app": rc.Name,
			},
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: &replicas,
			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{

				},
			},
		},
	}

	return rc.Interface.Create(crc)
}

func (rc *ReplicationController) Watch(completedHandler func() error) error {
	w, err := rc.Interface.Watch(metav1.ListOptions{
		LabelSelector: "app=" + rc.Name,
		Watch: true,
	})

	if err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return fmt.Errorf("failed on RC watching %v", ok)
			}
			switch event.Type {
			case watch.Added:
				//log.Info("Added: ", event.Object)
				//o := event.Object
				rc := event.Object.(*corev1.ReplicationController)
				log.Debug(rc.Name)
			case watch.Modified:
				rc := event.Object.(*corev1.ReplicationController)
				log.Debugf("RC: %s, Replicas: %d, AvailableReplicas: %d", rc.Name, rc.Status.Replicas, rc.Status.AvailableReplicas)
				if rc.Status.Replicas != 0 && rc.Status.AvailableReplicas == rc.Status.Replicas {
					var err error
					if nil !=  completedHandler {
						err = completedHandler()
					}
					w.Stop()
					return err
				}

			case watch.Deleted:
				log.Info("Deleted: ", event.Object)
			default:
				log.Error("Failed")
			}
		}
	}
}
