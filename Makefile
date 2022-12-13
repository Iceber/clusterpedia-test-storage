STORAGE_PLUGIN ?= sample-storage-layer

REGISTRY ?= "ghcr.io/clusterpedia-io/clusterpedia"
CLUSTERPEDIA_BUILDER_IMAGE = "ghcr.io/clusterpedia-io/clusterpedia/builder"
CLUSTERPEDIA_VERSIONS = v0.6.0-beta.1
RELEASE_ARCHS ?= amd64 arm64

BUILDER_IMAGE ?= ""

VERSION = $(shell git describe --tags 2>/dev/null)
ifeq ($(VERSION),)
	VERSION = v0.0.0
endif

BUILDER_TAG ?= $(shell echo $(BUILDER_IMAGE)|awk -F ':' '{ print $$2 }')
ifeq ($(BUILDER_TAG),)
	BUILDER_TAG = latest
endif

GOARCH ?= $(shell go env GOARCH)

PWD = $(shell pwd)
CLUSTERPEDIA_REPO ?= $(PWD)/clusterpedia

build-plugin:
	CLUSTERPEDIA_REPO=$(CLUSTERPEDIA_REPO) \
		clusterpedia/hack/builder.sh plugins $(STORAGE_PLUGIN).so

build-components:
	OUTPUT_DIR=$(PWD) ON_PLUGINS=true \
		$(MAKE) -C clusterpedia all

image-plugin:
ifeq ($(BUILDER_IMAGE), "")
	$(error BUILDER_IMAGE is not define)
endif

	docker buildx build \
		-t $(REGISTRY)/$(STORAGE_PLUGIN)-$(GOARCH):$(VERSION)-$(BUILDER_TAG) \
		--platform=linux/$(GOARCH) \
		--load \
		--build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--build-arg PLUGIN_NAME=$(STORAGE_PLUGIN).so .

push-images: clean-manifests
	set -e; \
	for version in $(CLUSTERPEDIA_VERSIONS); do \
	    images=""; \
	    for arch in $(RELEASE_ARCHS); do \
			GOARCH=$$arch BUILDER_IMAGE=$(CLUSTERPEDIA_BUILDER_IMAGE):$$version BUILDER_TAG=$$version $(MAKE) image-plugin; \
			image=$(REGISTRY)/$(STORAGE_PLUGIN)-$$arch:$(VERSION)-$$version; \
			docker push $$image; \
			images="$$images $$image"; \
		done; \
		docker manifest create $(REGISTRY)/$(STORAGE_PLUGIN):$(VERSION)-$$version $$images; \
		docker manifest push $(REGISTRY)/$(STORAGE_PLUGIN):$(VERSION)-$$version; \
	done;

clean-manifests:
	for version in $(CLUSTERPEDIA_VERSIONS); do \
		docker manifest rm $(REGISTRY)/$(STORAGE_PLUGIN):$(VERSION)-$$version 2>/dev/null; \
	done; exit 0
	
