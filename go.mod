module github.com/mvallim/gluster-simple-csi-driver

go 1.13

require (
	github.com/container-storage-interface/spec v1.3.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.7.0
	github.com/kubernetes-csi/drivers v1.0.2
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	google.golang.org/grpc v1.29.1
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.1-beta.0
	k8s.io/client-go v0.17.0
	k8s.io/component-base v0.17.0
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
)
