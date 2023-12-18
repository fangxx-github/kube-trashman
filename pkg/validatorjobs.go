package pkg

import (
	"fmt"
	_ "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/batch/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/kubernetes"
	"strconv"
	_ "strconv"
	"time"
)

type JobValidator func(v1.Job) (bool, error)

func ExpiredJobs(pendingJobExpirationtimes float64, annotation string) JobValidator {
	//识别job中过期的注释kube.trashman.io/expiration
	expirationAnnotationName := fmt.Sprintf("%s/expiration", annotation)

	return func(job v1.Job) (bool, error) {
		var age time.Duration
		if job.Status.CompletionTime == nil {
			for _, condition := range job.Status.Conditions {
				if condition.Reason == "BackoffLimitExceeded" && condition.Status == "True" {
					age = time.Since(condition.LastProbeTime.Time)
					break
				}
			}
		} else {
			age = time.Since(job.Status.CompletionTime.Time)
		}
		//判断如果间隔时间等于0 说明是刚创建的，不算过期
		if age == 0 {
			return false, nil
		}

		//识别job中过期的注释kube.trashman.io/expiration并复制给expirationAnnotationtime转化成浮点数
		expirationAnnotationtime := job.ObjectMeta.Annotations[expirationAnnotationName]
		if expirationAnnotationtime == "" {
			expirationAnnotationtime = "120.0"
		}
		maxAgeOverride, err := strconv.ParseFloat(expirationAnnotationtime, 64)
		if err != nil {
			log.Fatal(err.Error())
			return false, nil
		} else {
			log.Debugf(
				"Expiration override for (%s:%s) with annotation (%s) of (%s)",
				job.ObjectMeta.Namespace,
				job.ObjectMeta.Name,
				expirationAnnotationName,
				expirationAnnotationtime,
			)
			if age.Minutes() >= maxAgeOverride {
				return true, nil
			}
		}
		return age.Minutes() >= pendingJobExpirationtimes, nil
	}
}
