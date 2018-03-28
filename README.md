Lets break stuff and eat cookies
================================

CookieMonster will eat your applications. And infrastructure. And possibly your cat if it gets too close.

### Pre-requirements
 * Go lang 
 * Docker
 * Kubernetes
```sh 
export GOPATH=$HOME/go
git clone http://github.com/sktelecom-oslab/cookiemonster $GOPATH/src/cookiemonster
cd $GOPATH/src/cookiemonster
```

### Build
```sh
make vendor
make build-linux
make docker
```

### Run locally without Docker

##### Mac
```
./bin/cookiemonster-darwin-amd64 &
```
or
##### Linux
```
./bin/cookiemonster-linux-amd64 &
```

### Run locally with Docker
```
docker run -d -p 8080:8080 cookiemonster:latest
```

### Run on a Kubernetes cluster
```
helm install ./helm/cookiemonster-chart --name cookiemonster
```

Choose a host from the cluster and note the port that gets mapped
```
kubectl get svc cookiemonster
```

### Test it

##### Deploy stuff
```
kubectl create ns test
kubectl create clusterrolebinding test --clusterrole=cluster-admin --serviceaccount=test:default | true
kubectl create clusterrolebinding test --clusterrole=cluster-admin --serviceaccount=default:default | true
helm install skt/etcd --name etcdtest --namespace test --version 0.1.0
helm install skt/rabbitmq --name rabbitmqtest --namespace test --set replicas=7 --version 0.1.0
helm install skt/mariadb --name mariadbtest --namespace test --set replicas=7 --version 0.1.0
```

### URL endpoints to test
```
URL="localhost:30003"

curl http://$URL/list
curl http://$URL/list/killpod
curl http://$URL/list/nodeexec
curl http://$URL/show/killpod/openstack-random-deployment
curl http://$URL/show/nodeexec/test-exec-random

curl http://$URL/status/running
curl http://$URL/status/completed

curl http://$URL/start/killpod/openstack-random-deployment
curl http://$URL/stop/killpod/openstack-random-deployment

curl http://$URL/start/nodeexec/test-exec-random
curl http://$URL/stop/nodeexec/test-exec-random
```
