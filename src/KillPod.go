package main

import (
	"log"
	"math/rand"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

var server string = "https://10.0.1.17:6443"
var token string = "oreo0430@"
var caData []byte = []byte(
	`-----BEGIN CERTIFICATE-----
MIIC9zCCAd+gAwIBAgIJAOYCnu51/YB3MA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV
BAMMB2t1YmUtY2EwHhcNMTcwNDE4MDUzMTE4WhcNNDQwOTAzMDUzMTE4WjASMRAw
DgYDVQQDDAdrdWJlLWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
y+i0luxHveyXMRqvLlW56npZIkm3FynV9/Jd6et04pXtTpIOOwo8it9oPU4npOQv
NRbpUdCnoPzQrqXO3pXQAAiL2iZ1Nsg3mrBk+6wj9WT85XU2lR74eR//n/pjEtwd
3rFoauVUTBdKYd+orksAp+5RwFd1nhXxOpVAmCq35JZp1yI3GXMljnfMQ3ygQ1XS
nIbh0IgshgFqPd4iSBiJ4IADfQEB1ZlgoCu6IRlrwfNvaegO8XNIzskwtn7cVlRV
TNIQxZSvOuhTCxPMcT5RWgNHQO8VPDoZdLPnzwBJLVfi8ex429CisF2HmQtZV6pQ
SD0EQ+Xvpypz16YwpevmMQIDAQABo1AwTjAdBgNVHQ4EFgQUnZ55fF7dQgAwUsAg
N5ypWh/ee/kwHwYDVR0jBBgwFoAUnZ55fF7dQgAwUsAgN5ypWh/ee/kwDAYDVR0T
BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAXK/Y/nHoHMLB9a8iIezULcH2mxuo
uhzU/UmrojpSJcOfH6dtBsdktJmtnz/ujNW/DeZ1TZskse+7rnX37DkkPGHwVqQ3
/ayM7qQZMCJP6e97yv5FuDX1iYcxPTqJHDtrGF6dZBtxOexvlyXoBaApA2+gIxR6
IpR/2o/xxMSpXf/i2dxTwOLoSgCpn9mVgcK9lPhS3aIKSUPX69Rgvajzu0svjNv7
7/gGrFQLOgEafRYyS5HryRW2UyclsptcqjdyUHagI31ItoaNGCDJ4ksIfHfGv3RE
nU9FCNAFG774IGrQVS44aDGkH+C7MHW6eM5hQ7XVGbd2CAhVgyyhuWuA9g==
-----END CERTIFICATE-----`)

func randomInt(i int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(i)
}

func kubeConnect() (*kubernetes.Clientset, error) {
	var config *rest.Config
	if inKubeCluster {
		var err error
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	} else {
		config = &rest.Config{
			Host:            server,
			BearerToken:     token,
			TLSClientConfig: rest.TLSClientConfig{CAData: caData},
		}
	}

	return kubernetes.NewForConfig(config)
}

func victimDeployment(c *kubernetes.Clientset, ns, name string) (appsv1beta1.Deployment, bool) {
	var deployment appsv1beta1.Deployment
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := v1.ListOptions{}
		deployments, err := c.AppsV1beta1().Deployments(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(deployments.Items))
		deployment = deployments.Items[x]
		found = true
	} else {
		// query named deployment set via name provided
		lo := v1.ListOptions{FieldSelector: "metadata.name=" + name}
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

func victimStatefulSet(c *kubernetes.Clientset, ns, name string) (appsv1beta1.StatefulSet, bool) {
	var statefulset appsv1beta1.StatefulSet
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := v1.ListOptions{}
		statefulsets, err := c.AppsV1beta1().StatefulSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(statefulsets.Items))
		statefulset = statefulsets.Items[x]
		found = true
	} else {
		// query named statefulset via name provided
		lo := v1.ListOptions{FieldSelector: "metadata.name=" + name}
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

func victimDaemonSet(c *kubernetes.Clientset, ns, name string) (extv1beta1.DaemonSet, bool) {
	var daemonset extv1beta1.DaemonSet
	var found bool

	if name == "" {
		// no name specified, choose one at random
		lo := v1.ListOptions{}
		daemonsets, err := c.ExtensionsV1beta1().DaemonSets(ns).List(lo)
		if err != nil {
			panic(err.Error())
		}
		x := randomInt(len(daemonsets.Items))
		daemonset = daemonsets.Items[x]
		found = true
	} else {
		// query named daemonset via name provided
		lo := v1.ListOptions{FieldSelector: "metadata.name=" + name}
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

// choose a pod of parent 'kind' from 'namespace' and kind 'n' of them
func killPod(queryName, kind, ns string, n int, slackOut bool) string {

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
	lo := v1.ListOptions{LabelSelector: s}

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
	c.CoreV1().Pods(ns).Delete(victimName, &v1.DeleteOptions{})

	return victimName
}
