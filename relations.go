package main

import (
	"reflect"
	"strings"
)

type Relation struct {
	Namespace  string `json:"namespace"`
	ParentName string `json:"parent_name"`
	ParentType string `json:"parent_type"`
	ChildName  string `json:"child_name"`
	ChildType  string `json:"child_type"`
}

func getRelation(r *[]Relation, p K8sParentResourceType, c K8sChildResourceType) {
	if c.MatchesSelectors(p.GetSelectors()) {
		*r = append(*r, Relation{
			Namespace:  c.GetNamespace(),
			ParentName: p.GetName(),
			ParentType: "k8s_" + strings.ToLower(strings.Split(reflect.TypeOf(p).String(), ".")[1]),
			ChildName:  c.GetName(),
			ChildType:  "k8s_" + strings.ToLower(strings.Split(reflect.TypeOf(c).String(), ".")[1]),
		})
	}
}

func getRelationList(namespaceResult K8sResources) []Relation {
	relationList := make([]Relation, 0)
	for _, pod := range namespaceResult.K8s_pods {
		for _, deployment := range namespaceResult.K8s_deployments {
			getRelation(&relationList, deployment, pod)
		}
		for _, daemonset := range namespaceResult.K8s_daemonsets {
			getRelation(&relationList, daemonset, pod)
		}
		for _, service := range namespaceResult.K8s_services {
			getRelation(&relationList, service, pod)
		}
		relationList = append(relationList, Relation{
			Namespace:  pod.Namespace,
			ParentName: pod.Node,
			ParentType: "k8s_node",
			ChildName:  pod.Name,
			ChildType:  "k8s_pod",
		})
	}
	for _, service := range namespaceResult.K8s_services {
		for _, deployment := range namespaceResult.K8s_deployments {
			getRelation(&relationList, service, deployment)
		}
		for _, daemonset := range namespaceResult.K8s_daemonsets {
			getRelation(&relationList, service, daemonset)
		}
	}
	return relationList
}
