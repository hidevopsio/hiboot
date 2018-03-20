package openshift

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestImageStreamCrud(t *testing.T) {
	imageStreamName := "is-test"
	namespace := "openshift"

	imageStream := &ImageStream{
		Name: imageStreamName,
		Namespace: namespace,
	}

	// create imagestream
	is, err := imageStream.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// get imagestream
	is, err = imageStream.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, imageStreamName, is.ObjectMeta.Name)

	// update imagestream

	// delete imagestream
	err = imageStream.Delete()
	assert.Equal(t, nil, err)
}


