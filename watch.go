package main

import (
	"fmt"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func dumpAndUpdate(namespace string, clientset *kubernetes.Clientset, resourceMap *ResourceMap) {
	k8sResources := getK8sResources(clientset, namespace)
	if !reflect.DeepEqual(k8sResources, resourceMap.Resources[namespace]) {
		resources := make(map[string]K8sResources)
		resources[namespace] = k8sResources
		rm := ResourceMap{
			Nodes:     resourceMap.Nodes,
			Resources: resources,
		}
		rm.dumpResources()
		resourceMap.Resources[namespace] = k8sResources
	}
}

func watch(clientset *kubernetes.Clientset, resourceMap *ResourceMap) {
	factory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			namespace := obj.(*v1.Pod).Namespace
			dumpAndUpdate(namespace, clientset, resourceMap)
		},
		DeleteFunc: func(obj interface{}) {
			namespace := obj.(*v1.Pod).Namespace
			dumpAndUpdate(namespace, clientset, resourceMap)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			namespace := newObj.(*v1.Pod).Namespace
			dumpAndUpdate(namespace, clientset, resourceMap)
		},
	})

	nodeInformer := factory.Core().V1().Nodes().Informer()
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			resourceMap.update(clientset)
		},
		DeleteFunc: func(obj interface{}) {
			resourceMap.update(clientset)
		},
	})

	// Start the informers
	factory.Start(nil)

	// Wait for the informers to sync
	if !cache.WaitForCacheSync(nil, podInformer.HasSynced) {
		fmt.Printf("Timed out waiting for caches to sync\n")
		return
	}
	select {}
}
