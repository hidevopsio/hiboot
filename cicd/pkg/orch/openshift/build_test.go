package openshift

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
	//"sync"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"
)

func init() {
	log.SetLevel("debug")
}

func TestBuildCrud(t *testing.T) {
	namespace := "demo-dev"
	appName := "hello-world"
	imageTag := "latest"
	s2iImageStream := "s2i-java:1.0.5"
	s2iImageStreamNamespace := "openshift"
	buildCmd := "mvn clean package -DskipTests -Djava.net.preferIPv4Stack=true"
	gitUrl := "https://github.com/john-deng/hello-world.git"
	mvnMirrorUrl := os.Getenv("MAVEN_MIRROR_URL")

	log.Debug(mvnMirrorUrl)

	var from corev1.ObjectReference


	buildV1Client, err := buildv1.NewForConfig(k8s.Config)
	assert.Equal(t, nil, err)

	buildConfigs := buildV1Client.BuildConfigs(namespace)

	log.Debugf("workDir: %v", os.Getenv("PWD"))

	imageStream := &ImageStream{
		Name:      appName,
		Namespace: namespace,
	}

	// create imagestream
	_, err = imageStream.Get()
	if errors.IsNotFound(err) {
		is, err := imageStream.Create()
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, is.ObjectMeta.Name)

		from = corev1.ObjectReference{
			Kind:      "ImageStreamTag",
			Name:      s2iImageStream,
			Namespace: s2iImageStreamNamespace,
		}
	} else {

		from = corev1.ObjectReference{
			Kind: "ImageStreamTag",
			Name: appName + ":" + imageTag,
			Namespace: namespace,
		}
	}

	buildConfigSPec := &v1.BuildConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app": appName,
			},
		},
		Spec: v1.BuildConfigSpec{

			// The runPolicy field controls whether builds created from this build configuration can be run simultaneously.
			// The default value is Serial, which means new builds will run sequentially, not simultaneously.
			RunPolicy: v1.BuildRunPolicy("Serial"),
			CommonSpec: v1.CommonSpec{

				Source: v1.BuildSource{
					Type:   v1.BuildSourceType(v1.BuildSourceGit),
					Git: &v1.GitBuildSource{
						URI: gitUrl,
					},
				},
				Strategy: v1.BuildStrategy{
					Type: v1.BuildStrategyType(v1.SourceBuildStrategyType),
					SourceStrategy: &v1.SourceBuildStrategy{
						From: from,
					},
				},
				Output: v1.BuildOutput{
					To: &corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: appName + ":" + imageTag,
					},
				},
			},
		},
	}

	bc, err := buildConfigs.Get(appName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		bc, err := buildConfigs.Create(buildConfigSPec)
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, bc.Name)
	} else {
		bc.Spec.CommonSpec.Strategy.SourceStrategy.From = from
		bc, err := buildConfigs.Update(bc)
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, bc.Name)
	}

	// Get build config
	bc, err = buildConfigs.Get(appName, metav1.GetOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// trigger build manually
	buildRequestCauses := []v1.BuildTriggerCause{}
	incremental := false
	buildTriggerCauseManualMsg := "Manually triggered"
	buildRequest := v1.BuildRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildRequest",
			APIVersion: "build.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app": appName,
			},
		},
		TriggeredBy: append(buildRequestCauses,
			v1.BuildTriggerCause{
				Message: buildTriggerCauseManualMsg,
			},
		),
		SourceStrategyOptions: &v1.SourceStrategyOptions{
			Incremental: &incremental,
		},
		Env: []corev1.EnvVar{
			{
				Name:  "MAVEN_MIRROR_URL",
				Value: mvnMirrorUrl, // for test only, it should be passed form client
			},
			{
				Name: "MAVEN_CLEAR_REPO",
				Value: "false",
			},
			{
				Name: "BUILD_CMD",
				Value: buildCmd,
			},
		},
		From: &from,
	}

	builds := buildV1Client.Builds(namespace)

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	b, err := buildConfigs.Instantiate(appName, &buildRequest)
	//	assert.Equal(t, nil, err)
	//	assert.Contains(t, b.Name, appName)
	//	for {
	//		bx, err := builds.Get(b.Name, metav1.GetOptions{})
	//		if err == nil && bx.Status.Phase == v1.BuildPhase(v1.BuildPhaseRunning) {
	//			wg.Done()
	//			log.Debugf("build %v is running...", bx.Name)
	//			break
	//		}
	//	}
	//}()
	//wg.Wait()

	b, err := buildConfigs.Instantiate(appName, &buildRequest)
	assert.Equal(t, nil, err)
	assert.Contains(t, b.Name, appName)
	for {
		time.Sleep(100 * time.Millisecond)
		bx, err := builds.Get(b.Name, metav1.GetOptions{})
		exit := false
		if err == nil {
			switch bx.Status.Phase {
			case v1.BuildPhase(v1.BuildPhaseComplete):
				log.Debugf("build %v is completed", bx.Name)
				exit = true
				break
			case v1.BuildPhase(v1.BuildPhaseFailed):
				log.Debugf("build %v is failed", bx.Name)
				exit = true
				break
			default:
				continue
			}

			if exit {
				break
			}
		}
	}

	log.Debug("Done")

	// Delete build config
	//err = buildConfigs.Delete(appName, &metav1.DeleteOptions{})
	//assert.Equal(t, nil, err)
	//
	//err = imageStream.Delete()
	//assert.Equal(t, nil, err)
}
