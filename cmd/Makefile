VERSION ?= v1.0.0
NAME=gluster-simple-csi-plugin

all: publish

publish: compile build push clean

compile:
	@echo "==> Building the project"
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME}

build:
	@echo "==> Building the docker image provisioner"
	@docker build -f Dockerfile.provisioner -t mvallim/gluster-csi-driver:$(VERSION)-provisioner .
	@echo "==> Building the docker image agent"
	@docker build -f Dockerfile.agent -t mvallim/gluster-csi-driver:$(VERSION)-agent .

push:
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-provisioner"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-provisioner
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-provisioner"
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)-agent"
	@docker push mvallim/gluster-csi-driver:$(VERSION)-agent
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)-agent"

clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: all push fetch build-image clean
