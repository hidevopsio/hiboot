package pipelines


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
	KubeApi KubApiT
)


type KubApiT struct {
	kubeconfig *string
	config *rest.Config
	clientSet *kubernetes.Clientset
}

func init() {

	var err error

	if os.Getenv("KUBE_CLIENT_MODE") == "external" {
		log.Info("Kubernetes External Client Mode")
		if home := homedir.HomeDir(); home != "" {
			KubeApi.kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			KubeApi.kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		KubeApi.config, err = clientcmd.BuildConfigFromFlags(os.Getenv("KUBE_MASTER_URL"), *KubeApi.kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Info("Kubernetes Internal Client Mode")
		KubeApi.config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientSet
	KubeApi.clientSet, err = kubernetes.NewForConfig(KubeApi.config)
	if err != nil {
		panic(err.Error())
	}
}

