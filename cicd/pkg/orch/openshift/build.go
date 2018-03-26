// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package openshift

import (
	"github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"k8s.io/apimachinery/pkg/api/errors"
)

type GitSource struct{
	Url string
	Ref string
	Secret string
}


type BuildConfig struct {
	AppName string
	Namespace string
	Git GitSource
	ImageTag string

	// use NewFrom when creating new buildConfig
	NewFrom corev1.ObjectReference
	From corev1.ObjectReference

	Client *buildv1.BuildV1Client // TODO do we need export it?
	BuildConfigs buildv1.BuildConfigInterface
	Builds buildv1.BuildInterface
}

// @Title NewBuildConfig
// @Description Create new BuildConfig Instance
// @Param namespace, appName, gitUrl, imageTag, s2iImageStream string
// @Return *BuildConfig, error
func NewBuildConfig(namespace, appName, gitUrl, gitRef, gitSecret, imageTag, s2iImageStream string) (*BuildConfig, error) {

	log.Debug("NewBuildConfig()")

	client, err := buildv1.NewForConfig(k8s.Config)
	buildConfig := &BuildConfig{
		Client:       client, // TODO do we need export it?
		BuildConfigs: client.BuildConfigs(namespace),
		Builds:       client.Builds(namespace),

		NewFrom: corev1.ObjectReference{
			Kind:      "ImageStreamTag",
			Name:      s2iImageStream,
			Namespace: "openshift",
		},

		From: corev1.ObjectReference{
			Kind:      "ImageStreamTag",
			Name:      appName + ":" + imageTag,
			Namespace: namespace,
		},

		AppName: appName,
		Namespace: namespace,
		Git: GitSource{
			Url: gitUrl,
			Ref: gitRef,
			Secret: gitSecret,
		},

		ImageTag: imageTag,

	}
	return buildConfig, err
}


// @Title Create
// @Description Create new BuildConfig
// @Param
// @Return *v1.BuildConfig, error
func (b *BuildConfig) Create() (*v1.BuildConfig, error) {
	log.Debug("BuildConfig.Create()")
	// create imagestream
	imageStream := &ImageStream{
		Name:      b.AppName,
		Namespace: b.Namespace,
	}

	var from corev1.ObjectReference
	_, err := imageStream.Get()
	if errors.IsNotFound(err) {
		_, err := imageStream.Create()
		if err != nil {
			return nil, err
		}
		from = b.NewFrom
	} else {
		from = b.From
	}

	// buildConfig
	buildConfig := &v1.BuildConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.AppName,
			Labels: map[string]string{
				"app": b.AppName,
			},
		},
		Spec: v1.BuildConfigSpec{

			// The runPolicy field controls whether builds created from this build configuration can be run simultaneously.
			// The default value is Serial, which means new builds will run sequentially, not simultaneously.
			RunPolicy: v1.BuildRunPolicy("Serial"),
			CommonSpec: v1.CommonSpec{

				Source: v1.BuildSource{
					Type: v1.BuildSourceType(v1.BuildSourceGit),
					Git: &v1.GitBuildSource{
						URI: b.Git.Url,
						Ref: b.Git.Ref,
					},
					SourceSecret: &corev1.LocalObjectReference{
						Name: b.Git.Secret,
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
						Name: b.AppName + ":" + b.ImageTag,
					},
				},
			},
		},
	}

	bc, err := b.BuildConfigs.Get(b.AppName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		bc, err = b.BuildConfigs.Create(buildConfig)
	} else {
		bc.Spec.CommonSpec.Strategy.SourceStrategy.From = from
		bc, err = b.BuildConfigs.Update(bc)
	}

	return bc, err
}


// @Title Get
// @Description Get BuildConfig
// @Param
// @Return *v1.BuildConfig, error
func (b *BuildConfig) Get() (*v1.BuildConfig, error) {
	log.Debug("BuildConfig.Get()")
	return b.BuildConfigs.Get(b.AppName, metav1.GetOptions{})
}

// @Title Delete
// @Description Delete BuildConfig
// @Param
// @Return error
func (b *BuildConfig) Delete() error {
	log.Debug("BuildConfig.Delet()")
	return b.BuildConfigs.Delete(b.AppName, &metav1.DeleteOptions{})
}


// @Title Build
// @Description Start build according to previous build config settings, it will produce new image build
// @Param repo string, buildCmd string
// @Return *v1.Build, error
func (b *BuildConfig) Build(repo string, buildCmd string) (*v1.Build, error) {
	log.Debug("BuildConfig.Build()")
	incremental := false
	buildTriggerCauseManualMsg := "Manually triggered"
	buildRequest := v1.BuildRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildRequest",
			APIVersion: "build.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: b.AppName,
			Labels: map[string]string{
				"app": b.AppName,
			},
		},
		TriggeredBy: append([]v1.BuildTriggerCause{},
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
				Value: repo, // for test only, it should be passed from the client
			},
			{
				Name:  "MAVEN_CLEAR_REPO",
				Value: "false",
			},
			{
				Name:  "BUILD_CMD",
				Value: buildCmd,
			},
		},
		From: &b.From,
	}

	return b.BuildConfigs.Instantiate(b.AppName, &buildRequest)
}


// @Title GetBuild
// @Description Get current build
// @Param
// @Return *v1.Build, error
func (b *BuildConfig) GetBuild() (*v1.Build, error) {
	log.Debug("BuildConfig.GetBuild()")
	return b.Builds.Get(b.AppName, metav1.GetOptions{})
}


// @Title GetBuildStatus
// @Description Get current build status
// @Param
// @Return v1.BuildPhase, error
func (b *BuildConfig) GetBuildStatus() (v1.BuildPhase, error) {
	log.Debug("BuildConfig.GetBuildStatus()")
	build, err := b.GetBuild()
	return build.Status.Phase, err
}
