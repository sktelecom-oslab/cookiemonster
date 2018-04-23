# Web service API

### /api/v1/config<name>
http action: GET
get config file


### /api/v1/pod/start<name>
http action: POST
Initiate abuser. Will return an error if an existing run is ongoing for that pod name

- namespace: namespaces in which pods will be deleted
- kind: deployment, daemonset, statetfulset
- name: kind's name
- target: how many pods to remove at a time
- interval: time in seconds between kills
- duration: length of time in seconds to run

Example
```
curl -X POST -d '{"namespace": [{ "name": "openstack", "resource": [{ "kind": "deployment", "name": null, "target": 1 }] }], "interval": 60, "duration": 600, "slack": true }' -H "Content-Type: application/json" 'http://localhost:8080/api/v1/pod/start'
```

### /api/v1/pod/stop<name>
http action: POST
Immediately stop a running job. Returns an error if that abuser does not exist
