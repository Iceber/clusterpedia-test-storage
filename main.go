package main

import (
	plugin "github.com/clusterpedia-io/sample-storage-layer/storage"
)

func init() {
	plugin.RegisterStorageLayer()
}
