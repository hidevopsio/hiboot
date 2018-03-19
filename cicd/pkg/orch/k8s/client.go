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

	if os.Getenv("KUBE_CLIENT_MODE") == "external" {
		log.Info("Kubernetes External Client Mode")
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		Config, err = clientcmd.BuildConfigFromFlags(os.Getenv("KUBE_MASTER_URL"), *kubeconfig)
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
