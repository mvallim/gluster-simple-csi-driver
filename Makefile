VERSION ?= v1.0.0
NAME=gluster-simple-csi-plugin

all: publish

publish: build push

build:
	@echo "==> Building the docker image provisioner"
	@docker build -f Dockerfile.provisioner --rm -t mvallim/gluster-csi-driver:$(VERSION)-provisioner .
	@echo "==> Building the docker image agent"
	@docker build -f Dockerfile.agent --rm -t mvallim/gluster-csi-driver:$(VERSION)-agent .

push:
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-provisioner"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-provisioner
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-provisioner"
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-agent"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-agent
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-agent"
