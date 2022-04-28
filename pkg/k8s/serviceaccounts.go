package k8s

import (
	"context"
	"github.com/dfds/kiam2irsa/pkg/logging"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	RoleArnAnnotationName      string = "eks.amazonaws.com/role-arn"
	RegionalStsAnnotationName  string = "eks.amazonaws.com/sts-regional-endpoints"
	RegionalStsAnnotationValue string = "true"
)

func GetSA(cmd *cobra.Command) {
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

	serviceAccounts, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, sa := range serviceAccounts.Items {
		name := sa.Name
		ns := sa.Namespace
		annotations := sa.Annotations
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
			sugar.Infof("Service account %s in namespace %s is not yet migrated to IRSA", name, ns)
		}
	}
}
