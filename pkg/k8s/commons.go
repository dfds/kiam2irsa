package k8s

import (
	"github.com/dfds/kiam2irsa/pkg/logging"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

const (
	RoleAnnotation             string = "iam.amazonaws.com/role"
	RoleArnAnnotationName      string = "eks.amazonaws.com/role-arn"
	RegionalStsAnnotationName  string = "eks.amazonaws.com/sts-regional-endpoints"
	RegionalStsAnnotationValue string = "true"
	Parallelism                bool   = false
)

var nsWaitGroup sync.WaitGroup
var podWaitGroup sync.WaitGroup

func k8sClientSet(cmd *cobra.Command) (*kubernetes.Clientset, error) {
	sugar := logging.SugarLogger()

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}
	return clientset, nil
}
