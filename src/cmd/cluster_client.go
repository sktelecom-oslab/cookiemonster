package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//lo := metav1.ListOptions{TypeMeta: metav1.TypeMeta{Kind : "Deployment"}}
	//lo := metav1.ListOptions{
	//	TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind : "statefulset"},
	//	FieldSelector: "status.phase=Running",
	//	}
	//geto := metav1.GetOptions{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind : "statefulset"}}
	deploymentsClient := clientset.AppsV1().Deployments("ceph")

	for {
		//pods, err := clientset.CoreV1().Pods("").List(lo)
		pods, err := deploymentsClient.List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		for _, p := range pods.Items {
			namespace := p.Namespace
			pod := p.Name
			_, err = deploymentsClient.Get(p.Name, metav1.GetOptions{})
			//_, err = clientset.CoreV1().Pods(namespace).Get(pod, geto)
			if errors.IsNotFound(err) {
				fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %s in namespace %s: %v\n",
					pod, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(10 * time.Second)

	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
