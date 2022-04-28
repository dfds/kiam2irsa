package k8s

import (
	"fmt"
	"github.com/dfds/kiam2irsa/pkg/logging"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetPods(cmd *cobra.Command) {
	sugar := logging.SugarLogger()

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := GetNamespacesWithPermittedAnnotation(clientset)
	if err != nil {
		panic(err.Error())
	}

	for _, ns := range namespaces {
		fmt.Println(ns)
	}

	/*pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		name := pod.Name
		ns := pod.Namespace
		annotations := pod.Annotations
		hasFavorable := false
		hasUndesirable := false
		for annoKey, annoValue := range annotations {
			if annoKey == RoleArnAnnotationName {
				hasFavorable = true
			}
			if annoKey == RegionalStsAnnotationName && annoValue == RegionalStsAnnotationValue {
				hasUndesirable = true
			}
		}
		if hasFavorable && !hasUndesirable {
			sugar.Infof("Pod %s in namespace %s is not yet migrated to IRSA", name, ns)
		}
	}
	*/
}
