package main

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	ClusterIP   string            `json:"cluster_ip"`
	ServiceType v1.ServiceType    `json:"service_type"`
	Labels      map[string]string `json:"labels"`
	Selectors   map[string]string `json:"selectors"`
}

func (s Service) GetName() string {
	return s.Name
}

func (s Service) GetNamespace() string {
	return s.Namespace
}

func (s Service) GetSelectors() map[string]string {
	return s.Selectors
}

// getServices retrieves the list of services in the specified namespace.
func getServices(clientset *kubernetes.Clientset, namespaceName string) ([]Service, error) {
	serviceList := make([]Service, 0)

	// Get the list of services
	services, err := clientset.CoreV1().Services(namespaceName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of services
	for _, service := range services.Items {

		serviceList = append(serviceList, Service{
			Name:        service.ObjectMeta.Name,
			Namespace:   service.ObjectMeta.Namespace,
			ClusterIP:   service.Spec.ClusterIP,
			ServiceType: service.Spec.Type,
			Labels:      service.Labels,
			Selectors:   service.Spec.Selector,
		})
	}

	return serviceList, nil
}
