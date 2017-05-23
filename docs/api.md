# Web service API

## killpod
kill one or more pods

### /killpod/status/<name>
http action: GET
get action status for the a give pod name


### /killpod/start/<name>
http action: POST
Initiate abuser. Will return an error if an existing run is ongoing for that pod name

- type: deployment, daemonset, pod
- targets: how many pods to remove at a time. default is 1
- interval: time in seconds between kills, 0 for single, default is 0
- duration: length of time in seconds to run, 0 for continuous, default is 0

Example
```
curl -X POST -d '{"kind": "deployment", "target": 1, "interval": 30, "duration": 600}' -H "Content-Type: application/json" 'http://localhost:8080/killpod/start/mysql'
```

### /killpod/stop/<name>
http action: POST
Immediately stop a running job. Returns an error if that abuser does not exist


killservice
killnode
