package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Daemonset struct {
	Namespace              string            `json:"namespace"`
	Name                   string            `json:"name"`
	NumberAvailable        int32             `json:"number_available"`
	NumberUnavailable      int32             `json:"number_unavailable"`
	DesiredNumberScheduled int32             `json:"desired_number_scheduled"`
	CurrentNumberScheduled int32             `json:"current_number_scheduled"`
	NumberReady            int32             `json:"number_ready"`
	NumberMisscheduled     int32             `json:"number_misscheduled"`
	UpdatedNumberScheduled int32             `json:"updated_number_scheduled"`
	Containers             map[string]string `json:"containers"`
	Labels                 map[string]string `json:"labels"`
	Selectors              map[string]string `json:"selectors"`
}

func (d Daemonset) GetName() string {
	return d.Name
}

func (d Daemonset) GetNamespace() string {
	return d.Namespace
}

func (d Daemonset) GetSelectors() map[string]string {
	return d.Selectors
}

func (d Daemonset) MatchesSelectors(selectors map[string]string) bool {
	for key, value := range selectors {
		if d.Labels[key] != value {
			return false
		}
	}
	return true
}

// getdaemonsets retrieves the list of daemonsets in the specified namespace.
func getDaemonsets(clientset *kubernetes.Clientset, namespaceName string) ([]Daemonset, error) {
	daemonsetList := make([]Daemonset, 0)

	// Get the list of daemonsets
	daemonsets, err := clientset.AppsV1().DaemonSets(namespaceName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of daemonsets
	for _, daemonset := range daemonsets.Items {

		containerMap := make(map[string]string)
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			containerMap[container.Name] = container.Image
		}
		daemonsetList = append(daemonsetList, Daemonset{
			Name:                   daemonset.ObjectMeta.Name,
			Namespace:              daemonset.ObjectMeta.Namespace,
			NumberAvailable:        daemonset.Status.NumberAvailable,
			NumberUnavailable:      daemonset.Status.NumberUnavailable,
			DesiredNumberScheduled: daemonset.Status.DesiredNumberScheduled,
			CurrentNumberScheduled: daemonset.Status.CurrentNumberScheduled,
			NumberReady:            daemonset.Status.NumberReady,
			NumberMisscheduled:     daemonset.Status.NumberMisscheduled,
			UpdatedNumberScheduled: daemonset.Status.UpdatedNumberScheduled,
			Containers:             containerMap,
			Labels:                 daemonset.Labels,
			Selectors:              daemonset.Spec.Selector.MatchLabels,
		})
	}

	return daemonsetList, nil
}
