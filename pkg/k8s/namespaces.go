package k8s

import (
	"context"
	"github.com/dfds/kiam2irsa/pkg/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"regexp"
)

func GetNamespaces(clientset *kubernetes.Clientset) (*v1.NamespaceList, error) {
	sugar := logging.SugarLogger()
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}
	return namespaces, nil
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
		for annoKey, annoVal := range annotations {
			matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/.*", []byte(annoVal))
			if annoKey == "iam.amazonaws.com/permitted" && matchVal {
				retNamespaces = append(retNamespaces, name)
			}
		}
	}

	return retNamespaces, nil
}
