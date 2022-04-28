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

	namespaces, err := GetNamespacesWithPermittedAnnotation(clientset)
	if err != nil {
		sugar.Panic(err.Error())
	}

	for _, ns := range namespaces {
		nsWaitGroup.Add(1)
		go checkNamespace(clientset, ns)
	}
	nsWaitGroup.Wait()
}

func checkNamespace(clientset *kubernetes.Clientset, namespace string) {
	defer nsWaitGroup.Done()
	sugar := logging.SugarLogger()
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
	}
	for _, pod := range pods.Items {
		podWaitGroup.Add(1)
		go checkPod(clientset, pod)
	}
	podWaitGroup.Wait()
}

func checkPod(clientset *kubernetes.Clientset, pod v1.Pod) {
	defer podWaitGroup.Done()
	sugar := logging.SugarLogger()
	podName := pod.Name
	ns := pod.Namespace
	podAnnotations := pod.Annotations
	hasKiamAnnotation := false
	hasServiceAccountName := false
	hasIrsaAnnotation := false

	for annoKey, annoValue := range podAnnotations {
		matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/", []byte(annoValue))
		if annoKey == RoleAnnotation && matchVal {
			hasKiamAnnotation = true
		}
	}

	serviceAccountName := pod.Spec.ServiceAccountName
	if serviceAccountName != "" {
		hasServiceAccountName = true
		hasIrsaAnnotation, _ = ServiceAccountHasAnnotationForIRSA(clientset, serviceAccountName, ns)
	}

	if (hasKiamAnnotation && !hasServiceAccountName) || (hasKiamAnnotation && hasServiceAccountName && !hasIrsaAnnotation) {
		sugar.Infof("Pod %s in namespace %s is not yet migrated to IRSA", podName, ns)
	}

	if hasKiamAnnotation && hasServiceAccountName && hasIrsaAnnotation {
		sugar.Infof("Pod %s in namespace %s is migrated to IRSA, but still supports KIAM", podName, ns)
	}

	if !hasKiamAnnotation && hasServiceAccountName && hasIrsaAnnotation {
		sugar.Infof("Pod %s in namespace %s is fully migrated to IRSA", podName, ns)
	}
}
