package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Deployment struct {
	Namespace           string            `json:"namespace"`
	Name                string            `json:"name"`
	AvailableReplicas   int32             `json:"available_replicas"`
	DesiredReplicas     int32             `json:"desired_replicas"`
	TotalReplicas       int32             `json:"total_replicas"`
	UnavailableReplicas int32             `json:"unavaible_replicas"`
	UpdatedReplicas     int32             `json:"updated_replicas"`
	Containers          map[string]string `json:"containers"`
	Labels              map[string]string `json:"labels"`
	Selectors           map[string]string `json:"selectors"`
}

func (d Deployment) GetName() string {
	return d.Name
}

func (d Deployment) GetNamespace() string {
	return d.Namespace
}

func (d Deployment) GetSelectors() map[string]string {
	return d.Selectors
}

func (d Deployment) MatchesSelectors(selectors map[string]string) bool {
	for key, value := range selectors {
		if d.Labels[key] != value {
			return false
		}
	}
	return true
}

// getDeployments retrieves the list of deployments in the specified namespace.
func getDeployments(clientset *kubernetes.Clientset, namespaceName string) ([]Deployment, error) {
	deploymentList := make([]Deployment, 0)

	// Get the list of deployments
	deployments, err := clientset.AppsV1().Deployments(namespaceName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of deployments
	for _, deployment := range deployments.Items {

		// Get the updated replicas using a retry loop
		var updatedReplicas int32
		err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			updatedDeployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.Background(), deployment.ObjectMeta.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			updatedReplicas = updatedDeployment.Status.UpdatedReplicas
			return nil
		})
		if err != nil {
			fmt.Printf("Failed to get updated replicas for deployment %s in namespace %s: %v\n", deployment.ObjectMeta.Name, namespaceName, err)
			updatedReplicas = deployment.Status.UpdatedReplicas
		}

		containerMap := make(map[string]string)
		for _, container := range deployment.Spec.Template.Spec.Containers {
			containerMap[container.Name] = container.Image
		}

		deploymentList = append(deploymentList, Deployment{
			Name:                deployment.ObjectMeta.Name,
			Namespace:           deployment.ObjectMeta.Namespace,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			DesiredReplicas:     deployment.Status.Replicas,
			TotalReplicas:       deployment.Status.Replicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
			UpdatedReplicas:     updatedReplicas,
			Containers:          containerMap,
			Labels:              deployment.Labels,
			Selectors:           deployment.Spec.Selector.MatchLabels,
		})
	}

	return deploymentList, nil
}
