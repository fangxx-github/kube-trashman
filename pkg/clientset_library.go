package pkg

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "os"
	"path/filepath"
)

func Clientset_bendi() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 使用kubeconfig中的当前上下文,加载配置文件
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 本地创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return clientset

}

func Clientset_incluster() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	// 集群内创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return clientset
}
