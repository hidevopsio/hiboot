package openshift

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

func init() {
	log.SetLevel("debug")
}

func TestBuildCrud(t *testing.T) {
	// put below configs in yaml file
	namespace := "demo-dev"
	appName := "hello-world"
	gitUrl := "http://gitlab.vpclub:8022/moses-demos/hello-world.git"
	gitRef := "master"
	imageTag := "latest"
	s2iImageStream := "s2i-java:1.0.5"
	buildCmd := "mvn clean package -Dmaven.test.skip=true -Djava.net.preferIPv4Stack=true"
	mvnMirrorUrl := os.Getenv("MAVEN_MIRROR_URL")
	log.Debug(mvnMirrorUrl)

	log.Debugf("workDir: %v", os.Getenv("PWD"))

	buildConfig, err := NewBuildConfig(namespace, appName, gitUrl, gitRef, imageTag, s2iImageStream)
	assert.Equal(t, nil, err)

	bc, err := buildConfig.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// Get build config
	bc, err = buildConfig.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// Build image stream
	build, err := buildConfig.Build(mvnMirrorUrl, buildCmd)
	assert.Equal(t, nil, err)
	assert.Contains(t, build.Name, appName)

	log.Debug("Done")
}


