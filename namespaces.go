package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getNamespaceList retrieves the list of namespaces in the cluster.
func getNamespaceList(clientset *kubernetes.Clientset) ([]string, error) {
	namespaceList := []string{}

	// Get the list of namespaces
	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of namespaces
	for _, ns := range nsList.Items {
		namespaceList = append(namespaceList, ns.ObjectMeta.Name)
	}

	return namespaceList, nil
}
