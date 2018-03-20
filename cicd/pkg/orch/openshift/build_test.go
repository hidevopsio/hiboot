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
	"sync"
)

func init()  {
	log.SetLevel("debug")
}

func TestBuildCrud(t *testing.T)  {
	namespace := "demo-dev"
	appName := "bc-test"
	imageTag := "latest"
	s2iImageStream := "s2i-java:latest"
	s2iImageStreamNamespace := "openshift"
	//gitToken := "qxgjgsy2GsHq1Qb5rxTo"

	buildV1Client, err := buildv1.NewForConfig(k8s.Config)
	assert.Equal(t, nil, err)

	buildConfigs := buildV1Client.BuildConfigs(namespace)

	log.Debugf("workDir: %v", os.Getenv("PWD"))

	imageStream := &ImageStream{
		Name: appName,
		Namespace: namespace,
	}

	// create imagestream
	is, err := imageStream.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, is.ObjectMeta.Name)

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
					Type: v1.BuildSourceType(v1.BuildSourceBinary),
					Binary: &v1.BinaryBuildSource{},
				},
				Strategy: v1.BuildStrategy{
					Type: v1.BuildStrategyType(v1.SourceBuildStrategyType),
					SourceStrategy: &v1.SourceBuildStrategy{
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      s2iImageStream,
							Namespace: s2iImageStreamNamespace,
						},
					},
				},
				Output: v1.BuildOutput{
					To: &corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: appName + ":" + imageTag,
					},
				},
			},
			//Triggers: []v1.BuildTriggerPolicy{
			//	{
			//		Type: v1.BuildTriggerType(v1.GitLabWebHookBuildTriggerType),
			//		GenericWebHook: &v1.WebHookTrigger{
			//			Secret: gitToken,
			//		},
			//	},
			//},
		},
	}
	bc, err := buildConfigs.Create(buildConfigSPec)
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// Get build config
	bc, err = buildConfigs.Get(appName, metav1.GetOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// trigger build manually
	buildRequestCauses := []v1.BuildTriggerCause{}
	incremental := true
	buildTriggerCauseManualMsg := "Manually triggered"
	buildRequest := v1.BuildRequest{
		TypeMeta: metav1.TypeMeta{
			Kind: "BuildRequest",
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
		From: &corev1.ObjectReference{},
		Binary: &v1.BinaryBuildSource{},
	}

	builds := buildV1Client.Builds(namespace)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		b, err := buildConfigs.Instantiate(appName, &buildRequest)
		assert.Equal(t, nil, err)
		assert.Contains(t, b.Name, appName)
		for {
			bx, err := builds.Get(b.Name, metav1.GetOptions{})
			if err == nil && bx.Status.Phase == v1.BuildPhase(v1.BuildPhaseRunning) {
				wg.Done()
				log.Debugf("build %v is running...", bx.Name)
				break
			}
		}
	}()
	wg.Wait()

	// Delete build config
	err = buildConfigs.Delete(appName, &metav1.DeleteOptions{})
	assert.Equal(t, nil, err)

	err = imageStream.Delete()
	assert.Equal(t, nil, err)
}
