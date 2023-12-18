package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"kube-trashman/pkg"
	"os"
	"path/filepath"
)

func main() {
	annotation := flag.String(
		"at",
		"kube.trashman.io",
		"Annotation prefix to check when deleting jobs",
	)
	expirationtime := flag.Float64(
		"et",
		120,
		"Expiration time on jobs (in minutes)",
	)
	namespace := flag.String(
		"ns",
		"",
		"Namespace to target when deleting jobs (by default all namespaces are targeted)",
	)
	/*pendingJobExpirationtime := flag.Float64(
		"pet",
		-1.0,
		"Set the time (in minutes) that jobs will be removed if they are still in the pending state.By default, jobs stuck in a pending state are not removed",
	)*/
	verbose := flag.Bool(
		"verbose",
		false,
		"Increase verbosity of logging",
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	//创建k8s连接
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
	clientsetbendi, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	//

	jobList, err := clientsetbendi.BatchV1().Jobs(*namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//定义即将删除的job切片
	validatorslice := []pkg.JobValidator{}
	validatorslice = append(validatorslice, pkg.ExpiredJobs(*expirationtime, *annotation))

	//输出最终将要删除的job列表
	targetJobs := pkg.TargetJobs(jobList.Items, *annotation, validatorslice)
	//定义删除策略
	deletePolicy := metav1.DeletePropagationForeground
	//
	for _, job := range targetJobs {
		log.Infof("Deleting (%s:%s)", job.ObjectMeta.Namespace, job.ObjectMeta.Name)
		err = clientsetbendi.BatchV1().Jobs(*namespace).Delete(context.Background(), job.ObjectMeta.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
		if err != nil {
			log.Warnf(err.Error())
		}
	}
}
