package main

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	ROLE_ARN_ANNOTATION_NAME      string = "eks.amazonaws.com/role-arn"
	REGIONAL_STS_ANNOTATION_NAME  string = "eks.amazonaws.com/sts-regional-endpoints"
	REGIONAL_STS_ANNOTATION_VALUE string = "true"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)
	sugar := logger.Sugar()

	home := homedir.HomeDir()

	// TODO: Replace this dirty hack with the cobra CLI
	kubeconfig := filepath.Join(home, ".kube", "config")
	args := os.Args
	if len(args) >= 2 {
		kubeconfig = args[1]
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
			if annoKey == ROLE_ARN_ANNOTATION_NAME {
				hasFavorable = true
			}
			if annoKey == REGIONAL_STS_ANNOTATION_NAME && annoValue == REGIONAL_STS_ANNOTATION_VALUE {
				hasUndesirable = true
			}
		}
		if hasFavorable && !hasUndesirable {
			sugar.Infof("Service account %s in namespace %s is not yet migrated to IRSA", name, ns)
		}
	}
}
