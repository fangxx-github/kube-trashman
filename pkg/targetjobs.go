package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/batch/v1"
	"strconv"
)

//再函数TargetJobs中主要是通过输入一个job列表，以及我们需要排除带有annotation注释的job，和我们JobValidator三个参数，输出
//一份job列表，这份列表将jobList.Items的job去除带有annotation注释的job以及JobValidator返回true的job组合成一份新的jobList

func TargetJobs(jobList []v1.Job, annotation string, validatorList []JobValidator) []v1.Job {
	targetjobslist := []v1.Job{}

	//排除带有ignore注释的job
	ignoreAnnotationName := fmt.Sprintf("%s/ignore", annotation)
	log.Debugf("(%d) jobs found", len(jobList))

	//遍历job查看是否有注释说明忽略job的删除
	for _, job := range jobList {
		ignoreAnnotation := job.ObjectMeta.Annotations[ignoreAnnotationName]
		//fmt.Println(ignoreAnnotation)
		if ignoreAnnotation == "" {
			ignoreAnnotation = "false"
		}
		//fmt.Println(ignoreAnnotation)
		ignore, err := strconv.ParseBool(ignoreAnnotation)
		if err == nil && ignore {
			log.Debugf(
				"Ignoring (%s:%s) with annotation (%s) of (%s)",
				job.ObjectMeta.Namespace,
				job.ObjectMeta.Name,
				ignoreAnnotationName,
				ignoreAnnotation,
			)
			continue
		}
		//判断job的删除字段是否为true
		for _, removeCheck := range validatorList {
			remove, err := removeCheck(job)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			if remove {
				targetjobslist = append(targetjobslist, job)
				break
			}
		}
	}
	log.Debugf("(%d) jobs to remove", len(targetjobslist))
	return targetjobslist
}
