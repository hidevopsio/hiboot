package openshift

import (
	"testing"
	"github.com/openshift/api/image/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

var (
	imageV1Client *imagev1.ImageV1Client
	imageStreams  imagev1.ImageStreamInterface
)

func init() {
	log.SetLevel("debug")
	var err error
	imageV1Client, err = imagev1.NewForConfig(k8s.Config)
	if err != nil {
		os.Exit(1)
	}
	imageStreams = imageV1Client.ImageStreams("openshift")
}

func TestImageStreamCrud(t *testing.T) {
	imageStreamName := "imagestream-test"
	namespace := "openshift"
	imageStreamSpec := &v1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: imageStreamName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": imageStreamName,
			},
		},
	}

	// create imagestream
	is, err := imageStreams.Create(imageStreamSpec)
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// get imagestream
	is, err = imageStreams.Get(imageStreamName, metav1.GetOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// update imagestream

	// delete imagestream
	err = imageStreams.Delete(imageStreamName, &metav1.DeleteOptions{})
	assert.Equal(t, nil, err)
}


