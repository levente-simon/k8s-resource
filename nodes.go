package main

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Node struct {
	Name              string            `json:"name"`
	Status            string            `json:"status"`
	CreationTimestamp string            `json:"creation_timestamp"`
	IP                []v1.NodeAddress  `json:"ip"`
	OS                string            `json:"os"`
	Kernel            string            `json:"kernel"`
	CPU               string            `json:"cpu"`
	Memory            string            `json:"memory"`
	Roles             []string          `json:"roles"`
	Labels            map[string]string `json:"labels"`
}

func (n Node) GetName() string {
	return n.Name
}

func (n Node) GetNamespace() string {
	return ""
}

func (n Node) MatchesSelectors(selectors map[string]string) bool {
	for key, value := range selectors {
		if n.Labels[key] != value {
			return false
		}
	}
	return true
}

// getPods retrieves the list of pods in the specified namespace.
func getNodes(clientset *kubernetes.Clientset) ([]Node, error) {
	nodeList := make([]Node, 0)

	// Get the list of pods
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through the list of pods
	for _, node := range nodes.Items {

		var roles []string
		for key := range node.Labels {
			if strings.HasPrefix(key, "node-role.kubernetes.io/") {
				roles = append(roles, key)
			}
		}

		nodeList = append(nodeList, Node{
			Name:              node.ObjectMeta.Name,
			Status:            string(node.Status.Phase),
			CreationTimestamp: node.ObjectMeta.CreationTimestamp.String(),
			IP:                node.Status.Addresses,
			OS:                node.Status.NodeInfo.OSImage,
			Kernel:            node.Status.NodeInfo.KernelVersion,
			CPU:               node.Status.Capacity.Cpu().String(),
			Memory:            node.Status.Allocatable.Memory().String(),
			Roles:             roles,
			Labels:            node.Labels,
		})
	}

	return nodeList, nil
}
