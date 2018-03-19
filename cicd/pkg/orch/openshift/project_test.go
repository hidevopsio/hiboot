package openshift

import (
	"testing"
	"github.com/openshift/api/project/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/hidevopsio/hi/cicd/pkg/orch/k8s"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

var (
	projectV1Client *projectv1.ProjectV1Client
	projects        projectv1.ProjectInterface
)

func init() {
	log.SetLevel("debug")
	var err error
	projectV1Client, err = projectv1.NewForConfig(k8s.Config)
	if err != nil {
		os.Exit(1)
	}
	projects = projectV1Client.Projects()
}

func TestProjectLit(t *testing.T) {
	// get all projects
	p, err := projects.List(metav1.ListOptions{})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(p.Items))
	log.Debugf("There are %d projects in the cluster", len(p.Items))
}

func TestProjectCrud(t *testing.T) {
	projectName := "project-crud"
	ps := &v1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
			Labels: map[string]string{
				"project": projectName,
			},
		},
	}
	var err error

	// create project
	ps, err = projects.Create(ps)
	assert.Equal(t, nil, err)

	// read project
	p, err := projects.Get(projectName, metav1.GetOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// TODO: update test is not passed yet
	// update project
	//np := &v1.Project{
	//	ObjectMeta: metav1.ObjectMeta{
	//		ResourceVersion: p.ObjectMeta.ResourceVersion,
	//		Name: projectName,
	//		Annotations: map[string]string{
	//			"openshift.io/display-name:": projectName + "-test",
	//		},
	//	},
	//}
	//p, err = projects.Update(np)
	//assert.Equal(t, nil, err)

	// delete project
	err = projects.Delete(projectName, &metav1.DeleteOptions{})
	assert.Equal(t, nil, err)

}

