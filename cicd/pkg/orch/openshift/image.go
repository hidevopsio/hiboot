package openshift


import (
	"github.com/openshift/api/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

type ImageStreamAPI interface {
	Create() (*v1.ImageStream, error)
	Get() (*v1.ImageStream, error)
	Delete() error
}

type ImageStream struct{
	Name string
	Namespace string
	Interface  imagev1.ImageStreamInterface
}

var (
	imageV1Client *imagev1.ImageV1Client
)

func init() {
	log.SetLevel("debug")
	var err error
	imageV1Client, err = imagev1.NewForConfig(k8s.Config)
	if err != nil {
		os.Exit(1)
	}
}

func (is *ImageStream) Create() (*v1.ImageStream, error) {
	imageStream := &v1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: is.Name,
			Namespace: is.Namespace,
			Labels: map[string]string{
				"app": is.Name,
			},
		},
	}

	is.Interface = imageV1Client.ImageStreams(is.Namespace)

	// create image stream
	return is.Interface.Create(imageStream)
}

func (is *ImageStream) Get() (*v1.ImageStream, error) {
	return is.Interface.Get(is.Name, metav1.GetOptions{})
}

func (is *ImageStream) Delete() error {
	return is.Interface.Delete(is.Name, &metav1.DeleteOptions{})
}