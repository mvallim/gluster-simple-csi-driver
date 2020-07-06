VERSION ?= v1.0.0
NAME=gluster-simple-csi-plugin

all: publish

publish: build push
	@echo "==> Publish packages on Go"
	@curl -X GET "https://proxy.golang.org/github.com/mvallim/gluster-simple-csi-driver/@v/master.info"
	@echo "\n==> Published packages on Go"

build:
	@echo "==> Building the docker image controller"
	@docker build -f Dockerfile.controller --rm -t mvallim/gluster-csi-driver:$(VERSION)-controller .
	@echo "==> Building the docker image agent"
	@docker build -f Dockerfile.agent --rm -t mvallim/gluster-csi-driver:$(VERSION)-agent .

push:
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-controller"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-controller
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-controller"
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-agent"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-agent
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-agent"
