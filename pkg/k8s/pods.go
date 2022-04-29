package k8s

import (
	"context"
	"github.com/dfds/kiam2irsa/pkg/logging"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"regexp"
)

func CheckPodsMigrationStatus(cmd *cobra.Command) {
	sugar := logging.SugarLogger()
	clientset, err := k8sClientSet(cmd)
	if err != nil {
		sugar.Panic(err.Error())
	}
	status, err := getStatusFlag(cmd)
	if err != nil {
		sugar.Error("Unable to get status flag. Setting to default value: KIAM")
		status = "KIAM"
	}
	parallelism, err := getParallelismFlag(cmd)
	if err != nil {
		sugar.Error("Unable to determine parallelism. Setting to default value: false")
		parallelism = false
	}

	namespaces, err := GetNamespacesWithPermittedAnnotation(clientset)
	if err != nil {
		sugar.Panic(err.Error())
	}

	for _, ns := range namespaces {
		if parallelism {
			nsWaitGroup.Add(1)
			go checkAllPodsInNamespace(clientset, ns, status, parallelism)
		} else {
			checkAllPodsInNamespace(clientset, ns, status, parallelism)
		}

	}
	if parallelism {
		nsWaitGroup.Wait()
	}
}

func checkAllPodsInNamespace(clientset *kubernetes.Clientset, namespace string, status string, parallelism bool) {
	if parallelism {
		defer nsWaitGroup.Done()
	}
	sugar := logging.SugarLogger()
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
	}
	for _, pod := range pods.Items {
		if parallelism {
			podWaitGroup.Add(1)
			go checkPod(clientset, pod, status, parallelism)
		} else {
			checkPod(clientset, pod, status, parallelism)
		}
	}
	if parallelism {
		podWaitGroup.Wait()
	}
}

func checkPod(clientset *kubernetes.Clientset, pod v1.Pod, status string, parallelism bool) {
	if parallelism {
		defer podWaitGroup.Done()
	}
	sugar := logging.SugarLogger()
	podName := pod.Name
	ns := pod.Namespace
	hasKiamAnnotation := hasKiamAnnotation(pod)
	hasServiceAccountName, serviceAccountName := hasServiceAccountName(pod)
	hasIrsaAnnotation := false
	if hasServiceAccountName {
		hasIrsaAnnotation, _ = ServiceAccountHasAnnotationForIRSA(clientset, serviceAccountName, ns)
	}

	if status == "KIAM" {
		if isPodUsingKiam(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			sugar.Infof("Pod %s in namespace %s is using only KIAM", podName, ns)
		}
		return
	}

	if status == "BOTH" {
		if isPodUsingBoth(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			sugar.Infof("Pod %s in namespace %s is migrated to IRSA, but still supports KIAM", podName, ns)
		}
		return
	}

	if status == "IRSA" {
		if isPodUsingIrsa(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			sugar.Infof("Pod %s in namespace %s is fully migrated to IRSA", podName, ns)
		}
		return
	}
}

func hasKiamAnnotation(pod v1.Pod) bool {
	podAnnotations := pod.Annotations
	for annoKey, annoValue := range podAnnotations {
		matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/", []byte(annoValue))
		if annoKey == RoleAnnotation && matchVal {
			return true
		}
	}
	return false
}

func hasServiceAccountName(pod v1.Pod) (bool, string) {
	serviceAccountName := pod.Spec.ServiceAccountName
	if serviceAccountName != "" {
		return true, serviceAccountName
	}
	return false, serviceAccountName
}

func isPodUsingKiam(hasKiamAnnotation bool, hasIrsaAnnotation bool, hasServiceAccountName bool) bool {
	if (hasKiamAnnotation && !hasServiceAccountName) || (hasKiamAnnotation && hasServiceAccountName && !hasIrsaAnnotation) {
		return true
	}
	return false
}

func isPodUsingIrsa(hasKiamAnnotation bool, hasIrsaAnnotation bool, hasServiceAccountName bool) bool {
	if !hasKiamAnnotation && hasServiceAccountName && hasIrsaAnnotation {
		return true
	}
	return false
}

func isPodUsingBoth(hasKiamAnnotation bool, hasIrsaAnnotation bool, hasServiceAccountName bool) bool {
	if hasKiamAnnotation && hasServiceAccountName && hasIrsaAnnotation {
		return true
	}
	return false
}
