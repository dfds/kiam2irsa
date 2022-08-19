package k8s

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/dfds/kiam2irsa/pkg/logging"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckPodsMigrationStatus(cmd *cobra.Command) {
	sugar := logging.SugarLogger()
	clientset, err := k8sClientSet(cmd)
	if err != nil {
		sugar.Panic(err.Error())
	}
	status, err := getFlag(cmd, "status")
	if err != nil {
		sugar.Error("Unable to get status flag. Setting to default value: KIAM")
		status = "KIAM"
	}

	outputFormat, err := getFlag(cmd, "output")
	if err != nil {
		sugar.Error("Unable to get output flag. Setting to default value: TEXT")
		status = "TEXT"
	}

	checkAllPods(clientset, status, outputFormat)
}

func checkAllPods(clientset *kubernetes.Clientset, status string, outputFormat string) {
	var outputText bool = false
	if outputFormat == "TEXT" {
		outputText = true
	}
	sugar := logging.SugarLogger()
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
	}

	serviceAccounts, err := GetAllServiceAccounts(clientset)
	if err != nil {
		sugar.Panic(err.Error())
	}

	var wg sync.WaitGroup
	checkPodReturns := make(chan checkPodReturn, len(pods.Items))

	for _, pod := range pods.Items {
		wg.Add(1)
		pod := pod
		go func() {
			defer wg.Done()
			checkPod(pod, status, serviceAccounts, outputText, checkPodReturns)
		}()
	}

	wg.Wait()
	close(checkPodReturns)

	namespaceCount := make(map[string]int)
	if outputFormat == "CSV" {
		for podReturn := range checkPodReturns {
			if val, ok := namespaceCount[podReturn.namespace]; ok {
				namespaceCount[podReturn.namespace] = val + 1
			} else {
				namespaceCount[podReturn.namespace] = 1
			}
		}

		fmt.Println("namespace,count")
		for k, v := range namespaceCount {
			fmt.Printf("%s,%d\n", k, v)
		}
	}
}

func checkPod(pod v1.Pod, status string, saList *v1.ServiceAccountList, outputText bool, returnChan chan checkPodReturn) {
	sugar := logging.SugarLogger()
	podName := pod.Name
	ns := pod.Namespace
	hasKiamAnnotation := hasKiamAnnotation(pod)
	hasServiceAccountName, serviceAccountName := hasServiceAccountName(pod)
	hasIrsaAnnotation := false

	// This could be optimized to query all SAs in one go
	if hasServiceAccountName {
		hasIrsaAnnotation, _ = HasServiceAccountAnnotationForIRSA(serviceAccountName, ns, saList)
	}

	if status == "KIAM" {
		if isPodUsingKiam(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			if outputText {
				sugar.Infof("Pod %s in namespace %s is using only KIAM", podName, ns)
			}
			returnChan <- checkPodReturn{
				namespace: ns,
			}
		}
		return
	}

	if status == "BOTH" {
		if isPodUsingBoth(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			if outputText {
				sugar.Infof("Pod %s in namespace %s is migrated to IRSA, but still supports KIAM", podName, ns)
			}
			returnChan <- checkPodReturn{
				namespace: ns,
			}
		}
		return
	}

	if status == "IRSA" {
		if isPodUsingIrsa(hasKiamAnnotation, hasIrsaAnnotation, hasServiceAccountName) {
			if outputText {
				sugar.Infof("Pod %s in namespace %s is fully migrated to IRSA", podName, ns)
			}
			returnChan <- checkPodReturn{
				namespace: ns,
			}
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

type checkPodReturn struct {
	namespace string
}
