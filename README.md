# thanos-testing-toolbox
Tools for testing Thanos components easily and locally.

# Openshift local cluster

```bash
crc start
oc login -u kubeadmin -p XXX https://api.crc.testing:6443
oc new-project thanos
```

The server is accessible via web console at:

`https://console-openshift-console.apps-crc.testing`

For network errors:

```bash
oc edit dns.operator default
# or 
crc config set nameserver 8.8.8.8
```

Increase disk size:

```bash 
crc config set disk-size 80
crc config set memory 16384
```

Add
    
```yaml
    upstreams:
    - type: SystemResolvConf
    - address: 8.8.8.8
      port: 53
      type: Network
```

# Deploy Thanos

```bash
MANIFESTS_DIR=minio
oc delete -f 'manifests/$MANIFESTS_DIR/*' go run main.go && oc apply -f 'manifests/$MANIFESTS_DIR/*'
```

# Using custom images

```bash
make docker
docker tag thanosbench quay.io/rh-ee-tmange/thanosbench
docker push quay.io/rh-ee-tmange/thanosbench
```