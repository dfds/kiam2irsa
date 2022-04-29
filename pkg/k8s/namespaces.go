package k8s

import (
	"context"
	"github.com/dfds/kiam2irsa/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"regexp"
)

func GetAllNamespaces(clientset *kubernetes.Clientset) ([]string, error) {
	sugar := logging.SugarLogger()
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	retNamespaces := make([]string, 0)
	for _, ns := range namespaces.Items {
		retNamespaces = append(retNamespaces, ns.Name)
	}
	return retNamespaces, nil
}

func GetNamespacesWithPermittedAnnotation(clientset *kubernetes.Clientset) ([]string, error) {
	sugar := logging.SugarLogger()
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	retNamespaces := make([]string, 0)

	for _, ns := range namespaces.Items {
		name := ns.Name
		annotations := ns.Annotations
		for annoKey, annoValue := range annotations {
			matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/.*", []byte(annoValue))
			if annoKey == "iam.amazonaws.com/permitted" && matchVal {
				retNamespaces = append(retNamespaces, name)
			}
		}
	}

	return retNamespaces, nil
}
