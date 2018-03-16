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
	Client KubeClient
)


type KubeClient struct {
	Kubeconfig *string
	Config     *rest.Config
	ClientSet  *kubernetes.Clientset
}

func init() {

	var err error

	if os.Getenv("KUBE_CLIENT_MODE") == "external" {
		log.Info("Kubernetes External Client Mode")
		if home := homedir.HomeDir(); home != "" {
			Client.Kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			Client.Kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		Client.Config, err = clientcmd.BuildConfigFromFlags(os.Getenv("KUBE_MASTER_URL"), *Client.Kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Info("Kubernetes Internal Client Mode")
		Client.Config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the ClientSet
	Client.ClientSet, err = kubernetes.NewForConfig(Client.Config)
	if err != nil {
		panic(err.Error())
	}
}
