package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

type ResourceMap struct {
	Nodes     []Node                  `json:"k8s_nodes"`
	Resources map[string]K8sResources `json:"resources"`
}

func (rm *ResourceMap) dumpResources() {
	jsonString, err := json.Marshal(rm)
	if err != nil {
		log.Fatalf("Failed to marshal result to JSON: %v\n", err)
	}
	fmt.Println(string(jsonString))
}

func (rm *ResourceMap) update(clientset *kubernetes.Clientset) {
	resources := make(map[string]K8sResources)
	// Get the list of namespaces
	namespaceList, err := getNamespaceList(clientset)
	if err != nil {
		log.Fatalf("Failed to get list of namespaces: %v", err)
	}

	// Iterate through the list of namespaces and get the pods, deployments, and services for each namespace
	for _, namespaceName := range namespaceList {
		resources[namespaceName] = getK8sResources(clientset, namespaceName)
	}
	nodes, err := getNodes(clientset)
	if err != nil {
		log.Fatalf("Failed to get list of nodes: %v", err)
	}

	rm.Nodes = nodes
	rm.Resources = resources
}

func main() {
	// Parse command line arguments
	var kubeconfig string
	var cfgWatch bool
	var cfgInCluster bool

	flag.BoolVar(&cfgInCluster, "in-cluster", false, "Use in-cluster config")
	flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "Path to the Kubernetes configuration file")
	flag.BoolVar(&cfgWatch, "w", false, "Continuously watch and report changes")
	flag.Parse()

	var resourceMap ResourceMap
	// Create a Kubernetes clientset
	clientset, err := createClientset(cfgInCluster, kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}
	resourceMap.update(clientset)

	if cfgWatch {
		resourceMap.dumpResources()
		watch(clientset, &resourceMap)
	} else {
		resourceMap.dumpResources()
	}
}
