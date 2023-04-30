package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
)

type K8sResources struct {
	K8s_pods        []Pod        `json:"k8s_pods"`
	K8s_deployments []Deployment `json:"k8s_deployments"`
	K8s_daemonsets  []Daemonset  `json:"k8s_damonsets"`
	K8s_services    []Service    `json:"k8s_services"`
	Relations       []Relation   `json:"relations"`
}

type K8sResourceType interface {
	GetName() string
	GetApp() string
	GetNamespace() string
}

type K8sParentResourceType interface {
	GetName() string
	GetNamespace() string
	GetSelectors() map[string]string
}

type K8sChildResourceType interface {
	GetName() string
	GetNamespace() string
	MatchesSelectors(selectors map[string]string) bool
}

func getK8sResources(clientset *kubernetes.Clientset, namespaceName string) K8sResources {
	podResult, err := getPods(clientset, namespaceName)
	if err != nil {
		log.Fatalf("Failed to get pods in namespace %s: %v\n", namespaceName, err)
	}

	deploymentResult, err := getDeployments(clientset, namespaceName)
	if err != nil {
		log.Fatalf("Failed to get deployments in namespace %s: %v\n", namespaceName, err)
	}

	daemonsetResult, err := getDaemonsets(clientset, namespaceName)
	if err != nil {
		log.Fatalf("Failed to get daemonsets in namespace %s: %v\n", namespaceName, err)
	}

	serviceResult, err := getServices(clientset, namespaceName)
	if err != nil {
		log.Fatalf("Failed to get services in namespace %s: %v\n", namespaceName, err)
	}

	namespaceResult := K8sResources{
		K8s_pods:        podResult,
		K8s_deployments: deploymentResult,
		K8s_daemonsets:  daemonsetResult,
		K8s_services:    serviceResult,
		Relations:       nil,
	}
	namespaceResult.Relations = getRelationList(namespaceResult)

	return namespaceResult
}
