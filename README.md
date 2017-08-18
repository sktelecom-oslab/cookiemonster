## Lets break stuff and eat cookies

CookieMonster will eat your applications. And infrastructure. And possibly your cat if it gets too close.

### Build
```
make vendor
make build
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
docker run -d -p 8080:8080 oreo01:5000/cookiemonster:latest
```

### Run on a Kubernetes cluster
```
kubectl create -f ./k8s/
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

### v2 ###
```
URL="localhost:8080"
# curl http://$URL/list
# curl http://$URL/list/killpod
# curl http://$URL/list/nodeexec
-- curl http://$URL/show/killpod/openstack-random-deployment
-- curl http://$URL/show/nodeexec/test-exec-random

curl http://$URL/status/running
curl http://$URL/status/completed

curl http://$URL/start/killpod/openstack-random-deployment
curl http://$URL/stop/killpod/openstack-random-deployment

curl http://$URL/start/nodeexec/test-exec-random
curl http://$URL/stop/nodeexec/test-exec-random
```



##### Running locally
```
curl -H "Content-Type: application/json" -X POST -d @./cookies.d/openstack-random-deployment.json 'http://localhost:8080/killpod/start/'
curl -H "Content-Type: application/json" -X POST -d @./cookies.d/openstack-random-deployment.json 'http://localhost:8080/killpod/stop/'
```

##### Running on Kubernetes
```
curl -H "Content-Type: application/json" -X POST -d @./cookies.d/test-rabbitmq.json 'http://oreo07:32568/killpod/start/'
curl -H "Content-Type: application/json" -X POST -d @./cookies.d/test-exec-random.json 'http://oreo07:32568/nodeexec/start/'
```
