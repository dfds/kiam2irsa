package k8s

import (
	"context"
	"github.com/dfds/kiam2irsa/pkg/logging"
	"regexp"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckAllServiceAccounts(cmd *cobra.Command) {
	sugar := logging.SugarLogger()

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		sugar.Panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sugar.Panic(err.Error())
	}

	serviceAccounts, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
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

func ServiceAccountHasAnnotationForIRSA(clientset *kubernetes.Clientset, name string, namespace string) (bool, error) {
	sugar := logging.SugarLogger()
	serviceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		sugar.Error(err.Error())
		return false, err
	}

	hasRoleArn := false
	hasRegionalSts := false

	for key, val := range serviceAccount.Annotations {
		matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/", []byte(val))
		if key == RoleArnAnnotationName && matchVal {
			hasRoleArn = true
		}
		if key == RegionalStsAnnotationName && val == RegionalStsAnnotationValue {
			hasRegionalSts = true
		}
	}

	if hasRoleArn && hasRegionalSts {
		return true, nil
	}
	return false, nil
}
