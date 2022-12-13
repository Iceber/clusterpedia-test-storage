# Sample Storage Layer
The storage layer plugin example replicates the implementation of the [internal storage layer](https://github.com/clusterpedia-io/clusterpedia/tree/main/pkg/storage/internalstorage), but it has the storage layer name - **sample storage layer**.

And you can see the storage layer information from the log when you list the resources.
```bash
$ STORAGE_PLUGINS=./plugins ./bin/apiserver --storage-name=sample-storage-layer --storage-config=./config.yaml <... other flags>
...

I1107 11:35:14.450738   18955 resource_storage.go:195] "list resources in the sample storage layer" gvr="pods"
I1107 11:35:14.483075   18955 httplog.go:131] "HTTP" verb="GET" URI="/apis/clusterpedia.io/v1beta1/resources/api/v1/pods?limit=1" latency="40.270678ms" userAgent="curl/7.64.1" audit-ID="4ed69024-1848-4aa1-9f72-a9f1d7f36e3f" srcIP="127.0.0.1:58568" resp=200
```
`curl -k --cert client.crt --key client.key https://127.0.0.1:8443/apis/clusterpedia.io/v1beta1/resources/api/v1/pods\?limit\=1`

## Build and Run
`git clone` repo
```bash
$ git clone --recursive https://github.com/clusterpedia-io/sample-storage-layer.git
$ cd sample-storage-layer
```

build storage layer plugin
```bash
$ make build-plugin

$ # check plugin
$ file ./plugins/sample-storage-layer.so
./plugins/sample-storage-layer.so: Mach-O 64-bit dynamically linked shared library x86_64
```

build clusterpedia components for the debug
```bash
$ make build-components
$ ls -al ./bin
drwxr-xr-x   6 icebergu  staff       192 11  7 11:17 .
drwxr-xr-x  16 icebergu  staff       512 11  7 11:15 ..
-rwxr-xr-x   1 icebergu  staff  90707488 11  7 11:15 apiserver
-rwxr-xr-x   1 icebergu  staff  91896016 11  7 11:16 binding-apiserver
-rwxr-xr-x   1 icebergu  staff  82769728 11  7 11:16 clustersynchro-manager
-rwxr-xr-x   1 icebergu  staff  45682000 11  7 11:17 controller-manager
```

run clusterpedia apiserver
```bash
$ STORAGE_PLUGINS=./plugins ./bin/apiserver --storage-name=sample-storage-layer --storage-config=./config.yaml <... other flags>
```

run clusterpedia clustersynchro-manager
```bash
$ STORAGE_PLUGINS=./plugins ./bin/clustersynchro-manager --storage-name=sample-storage-layer --storage-config=./config.yaml <... other flags>
```
