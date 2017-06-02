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
kubectl create clusterrolebinding test --clusterrole=cluster-admin --serviceaccount=test:default
kubectl create clusterrolebinding test --clusterrole=cluster-admin --serviceaccount=default:default
kubectl create secret generic ceph-secret-user --from-literal=key=AQCPrrdYbn6PJxAAWymornUUVbCkgWN6h+mPFg== --namespace=test
helm install skt/etcd --name etcdtest --namespace test --version 0.1.0
helm install skt/rabbitmq --name rabbitmqtest --namespace test --set replicas=7 --version 0.1.0
helm install skt/mariadb --name mariadbtest --namespace test --set replicas=7 --version 0.1.0
```

##### Running locally
```
curl -H "Content-Type: application/json" -X POST -d @./json/openstack-random-deployment-start.json 'http://localhost:8080/killpod/start/'
curl -H "Content-Type: application/json" -X POST -d @./json/openstack-random-deployment-stop.json 'http://localhost:8080/killpod/stop/'
```

##### Running on Kubernetes
```
curl -X POST -d '{"kind": "deployment", "target": 1, "interval": 30, "duration": 600}' -H "Content-Type: application/json" 'http://<host>:<port>/killpod/start/'
```
