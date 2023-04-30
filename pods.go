package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	Namespace         string            `json:"namespace"`
	Name              string            `json:"name"`
	Node              string            `json:"node"`
	Status            string            `json:"status"`
	CreationTimestamp string            `json:"creation_timestamp"`
	IP                string            `json:"ip"`
	Labels            map[string]string `json:"labels"`
}

func (p Pod) GetName() string {
	return p.Name
}

func (p Pod) GetNamespace() string {
	return p.Namespace
}

func (p Pod) MatchesSelectors(selectors map[string]string) bool {
	for key, value := range selectors {
		if p.Labels[key] != value {
			return false
		}
	}
	return true
}

// getPods retrieves the list of pods in the specified namespace.
func getPods(clientset *kubernetes.Clientset, namespaceName string) ([]Pod, error) {
	podList := make([]Pod, 0)

	// Get the list of pods
	pods, err := clientset.CoreV1().Pods(namespaceName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of pods
	for _, pod := range pods.Items {

		podList = append(podList, Pod{
			Name:              pod.ObjectMeta.Name,
			Namespace:         pod.ObjectMeta.Namespace,
			Node:              pod.Spec.NodeName,
			Status:            string(pod.Status.Phase),
			CreationTimestamp: pod.ObjectMeta.CreationTimestamp.Time.String(),
			IP:                pod.Status.PodIP,
			Labels:            pod.Labels,
		})
	}

	return podList, nil
}
