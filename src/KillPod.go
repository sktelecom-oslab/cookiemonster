package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var server string = "https://10.0.1.17:6443"
var token string = "oreo0430@"
var caData []byte = []byte(`-----BEGIN CERTIFICATE-----
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

// choose a pod of parent 'kind' from 'namespace' and kind 'n' of them
func killPod(kind string, namespace string, n int) {
	// config, err := rest.InClusterConfig()

	config := &rest.Config{
		Host:            server,
		BearerToken:     token,
		TLSClientConfig: rest.TLSClientConfig{CAData: caData},
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("could not connect to Kubernetes API: %s", err)
	}

	info, err := c.Discovery().ServerVersion()
	if err != nil {
		log.Fatal("failed to retrieve server API version: %s\n", err)
	}

	log.Printf("server API version information: %s\n", info)
}
