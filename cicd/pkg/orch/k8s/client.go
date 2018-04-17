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


package k8s


import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8s.io/client-go/rest"
	"os"
	"flag"
	"path/filepath"
	log "github.com/kataras/golog"
)

var (
	Config     *rest.Config
	ClientSet  *kubernetes.Clientset

	kubeconfig *string
)


func init() {

	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		log.Info("Kubernetes External Client Mode")
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		Config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Info("Kubernetes Internal Client Mode")
		Config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the ClientSet
	ClientSet, err = kubernetes.NewForConfig(Config)
	if err != nil {
		panic(err.Error())
	}
}
