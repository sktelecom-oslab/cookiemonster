package main

import (
	"log"
	"math/rand"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
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

func victimDeployment(c *kubernetes.Clientset, ns, name string) v1beta1.Deployment {

	// list existing deployments in the namespace
	d := v1.ListOptions{}
	deployments, err := c.AppsV1beta1().Deployments(ns).List(d)
	if err != nil {
		panic(err.Error())
	}

	var deployment v1beta1.Deployment

	if name == "" {
		// no name specified, choose one at random
		x := randomInt(len(deployments.Items))
		deployment = deployments.Items[x]
	} else {
		// find the specified deployment
		for _, d := range deployments.Items {
			if d.ObjectMeta.Name == name {
				deployment = d
			}
		}
	}
	return deployment
}

// choose a pod of parent 'kind' from 'namespace' and kind 'n' of them
func killPod(name, kind, ns string, n int, slackOut bool) string {
	var victimName string

	c, err := kubeConnect()
	if err != nil {
		panic(err.Error())
	}

	vd := victimDeployment(c, ns, name)
	if &vd == nil {
		log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
		return ""
	}

	log.Printf("Found %s %s in namespace %s\n", kind, vd.ObjectMeta.Name, ns)

	// convert Selector to ListOption
	s := ""
	for k, v := range vd.Spec.Selector.MatchLabels {
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
	victimName = pods.Items[x].ObjectMeta.Name
	log.Printf("Eating pod %s NOM NOM NOM!!!!", victimName)
	if slackOut {
		postSlack("Eating pod " + victimName + " from namespace " + ns + "!!!!  NOM! NOM! NOM!")
	}
	c.CoreV1().Pods(ns).Delete(victimName, &v1.DeleteOptions{})

	return victimName
}
