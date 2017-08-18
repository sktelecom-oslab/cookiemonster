package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clusterMode bool

func randomInt(i int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(i)
}

func kubeConnect() (*kubernetes.Clientset, error) {
	// attempt in kubebernetes cluster client init
	var config, err = rest.InClusterConfig()

	if err == nil {
		log.Println("Running in Kubernetes cluster")
		clusterMode = true
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return clientset, err
	} else {
		// attempt current context in kubeconfig
		var kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		log.Println("Running locally via ~/.kube/config")
		clusterMode = false
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return clientset, err
	}
}

func victimDeployment(c *kubernetes.Clientset, ns, name string) (apps_v1beta1.Deployment, bool) {
	var deployment apps_v1beta1.Deployment
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := meta_v1.ListOptions{}
		deployments, err := c.AppsV1beta1().Deployments(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(deployments.Items))
		deployment = deployments.Items[x]
		found = true
	} else {
		// query named deployment set via name provided
		lo := meta_v1.ListOptions{FieldSelector: "metadata.name=" + name}
		deployments, err := c.AppsV1beta1().Deployments(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		if len(deployments.Items) > 0 {
			deployment = deployments.Items[0]
			found = true
		}
	}

	return deployment, found
}

func victimStatefulSet(c *kubernetes.Clientset, ns, name string) (apps_v1beta1.StatefulSet, bool) {
	var statefulset apps_v1beta1.StatefulSet
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := meta_v1.ListOptions{}
		statefulsets, err := c.AppsV1beta1().StatefulSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(statefulsets.Items))
		statefulset = statefulsets.Items[x]
		found = true
	} else {
		// query named statefulset via name provided
		lo := meta_v1.ListOptions{FieldSelector: "metadata.name=" + name}
		statefulsets, err := c.AppsV1beta1().StatefulSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		if len(statefulsets.Items) > 0 {
			statefulset = statefulsets.Items[0]
			found = true
		}
	}

	return statefulset, found
}

func victimDaemonSet(c *kubernetes.Clientset, ns, name string) (ext_v1beta1.DaemonSet, bool) {
	var daemonset ext_v1beta1.DaemonSet
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := meta_v1.ListOptions{}
		daemonsets, err := c.ExtensionsV1beta1().DaemonSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(daemonsets.Items))
		daemonset = daemonsets.Items[x]
		found = true
	} else {
		// query named daemonset via name provided
		lo := meta_v1.ListOptions{FieldSelector: "metadata.name=" + name}
		daemonsets, err := c.ExtensionsV1beta1().DaemonSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		if len(daemonsets.Items) > 0 {
			daemonset = daemonsets.Items[0]
			found = true
		}
	}

	return daemonset, found
}

func victimNode(c *kubernetes.Clientset, name string) (*v1.Node, bool) {
	if name == "" {
		// no name specified, choose one at random
		lo := meta_v1.ListOptions{}
		nodes, err := c.CoreV1().Nodes().List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(nodes.Items))
		node := nodes.Items[x]
		return &node, true
	} else {
		// query named daemonset via name provided
		opt := meta_v1.GetOptions{}
		node, err := c.CoreV1().Nodes().Get(name, opt)
		if err != nil {
			panic(err.Error())
		}
		return node, true
	}
}

// choose a pod of parent 'kind' from 'namespace'
func killPod(queryName, kind, ns string, slackOut bool) string {

	c, err := kubeConnect()
	if err != nil {
		panic(err.Error())
	}

	// obtain the match labels from the specified 'kind'
	var matchLabels map[string]string

	switch kind {
	case "deployment":
		vd, found := victimDeployment(c, ns, queryName)
		name := vd.ObjectMeta.Name
		if !found {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return ""
		}
		// deployments should remain functional as long as a single pod is
		// available, but rabbitmq is an exception
		log.Printf("%s %s in namespace %s has %d pods defined, %d available and %d unavailable ",
			kind, name, ns, *vd.Spec.Replicas, vd.Status.AvailableReplicas, vd.Status.UnavailableReplicas)
		if name == "rabbitmq" {
			if (*vd.Spec.Replicas/2 + 1) >= vd.Status.AvailableReplicas {
				log.Printf("available pods less than unavailable pods, doing nothing")
				return ""
			}
		} else {
			if vd.Status.AvailableReplicas < 2 {
				log.Printf("Only one pod available, doing nothing")
				return ""
			}
		}

		log.Printf("Found %s %s in namespace %s\n", kind, name, ns)
		matchLabels = vd.Spec.Selector.MatchLabels
	case "statefulset":
		vss, found := victimStatefulSet(c, ns, queryName)
		name := vss.ObjectMeta.Name
		if !found {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return ""
		}
		log.Printf("%s %s in namespace %s has %d pods defined, %d available",
			kind, name, ns, *vss.Spec.Replicas, vss.Status.Replicas)
		if (*vss.Spec.Replicas/2 + 1) >= vss.Status.Replicas {
			log.Printf("available pods less than unavailable pods, doing nothing")
			return ""
		}

		log.Printf("Found %s %s in namespace %s\n", kind, vss.ObjectMeta.Name, ns)
		matchLabels = vss.Spec.Selector.MatchLabels
	case "daemonset":
		vds, found := victimDaemonSet(c, ns, queryName)
		name := vds.ObjectMeta.Name
		if !found {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return ""
		}
		log.Printf("%s %s in namespace %s has %d pods defined, %d available",
			kind, name, ns, vds.Status.DesiredNumberScheduled, vds.Status.NumberReady)
		if vds.Status.NumberReady < 2 {
			log.Printf("available pods less than unavailable pods, doing nothing")
			return ""
		}

		log.Printf("Found %s %s in namespace %s\n", kind, vds.ObjectMeta.Name, ns)
		matchLabels = vds.Spec.Selector.MatchLabels
	}

	// convert Selector to ListOption
	s := ""
	for k, v := range matchLabels {
		s = s + k + "=" + v + ","
	}
	s = s[:len(s)-1]
	lo := meta_v1.ListOptions{LabelSelector: s}

	// query pods
	pods, err := c.CoreV1().Pods(ns).List(lo)
	if err != nil {
		panic(err.Error())
	}

	x := randomInt(len(pods.Items))
	victimName := pods.Items[x].ObjectMeta.Name
	log.Printf("Eating pod %s NOM NOM NOM!!!!", victimName)
	if slackOut {
		postSlack("Eating pod " + victimName + " from namespace " + ns + "!!!!  NOM! NOM! NOM!")
	}
	c.CoreV1().Pods(ns).Delete(victimName, &meta_v1.DeleteOptions{})

	return victimName
}

func nodeIP(n *v1.Node) (string, error) {
	addresses := n.Status.Addresses
	for _, na := range addresses {
		if na.Type == v1.NodeInternalIP {
			return na.Address, nil
		}
	}
	return "", errors.New("Node " + n.ObjectMeta.Name + " IP not found")
}

func doExec(data ExecData) ([]byte, error) {
	c, err := kubeConnect()
	if err != nil {
		return nil, err
	}

	n, found := victimNode(c, data.Target)
	name := n.ObjectMeta.Name
	if !found {
		log.Printf("Can not find node %s", name)
		return nil, errors.New("AAAHHHH!!!")
	}
	ip, err := nodeIP(n)
	if err != nil {
		return nil, err
	}

	resp, err := sendCommands(ip, data.Commands)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
