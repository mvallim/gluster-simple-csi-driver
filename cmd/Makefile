VERSION ?= v0.0.0
NAME=gluster-csi-driver

all: publish

publish: compile build push clean

compile:
	@echo "==> Building the project"
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME}

build:
	@echo "==> Building the docker image"
	@docker build -t mvallim/gluster-csi-driver:$(VERSION) .

push:
	@echo "==> Publishing mvallim/gluster-csi-driver:$(VERSION)"
	@docker push mvallim/gluster-csi-driver:$(VERSION)
	@echo "==> Your image is now available at mvallim/gluster-csi-driver:$(VERSION)"

clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: all push fetch build-image clean