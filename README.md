## Lets break stuff and eat cookies

CookieMonster will eat your applications. And infrastructure. And possibly your cat if it gets too close.

### Build
```
make vendor
make build
make docker
```
- - -
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
- - -
### Test it

##### Running locally
```
curl -X POST -d '{"kind": "deployment", "target": 1, "interval": 30, "duration": 600}' -H "Content-Type: application/json" 'http://localhost:8080/killpod/start/mysql'
```

##### Running on Kubernetes
```
curl -X POST -d '{"kind": "deployment", "target": 1, "interval": 30, "duration": 600}' -H "Content-Type: application/json" 'http://<host>:<port>/killpod/start/mysql'
```
