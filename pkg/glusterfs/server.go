package glusterfs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"k8s.io/utils/mount"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

// NewDriver create new instance driver
func NewDriver(config *config.Config) (*Driver, error) {

	if config == nil {
		klog.Errorf("GlusterFS Simple CSI driver initialization failed")
		return nil, nil
	}

	driver := &Driver{
		Config: config,
	}

	klog.Infof("GlusterFS Simple CSI driver initialized")

	return driver, nil

}

// NewControllerServer create new instance controller
func NewControllerServer(driver *Driver) *ControllerServer {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()

	if err != nil {
		klog.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		klog.Fatalf("Failed to create client: %v", err)
	}

	restclient := clientset.CoreV1().RESTClient()

	return &ControllerServer{
		Driver:     driver,
		client:     clientset,
		config:     config,
		restclient: restclient,
	}
}

// NewNodeServer create new instance node
func NewNodeServer(driver *Driver) *NodeServer {
	return &NodeServer{
		Driver:  driver,
		mounter: mount.New(""),
	}
}

// NewIdentityServer create new instance identity
func NewIdentityServer(driver *Driver) *IdentityServer {
	return &IdentityServer{
		Driver: driver,
	}
}

// Run start server
func (driver *Driver) Run() {
	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(driver.Endpoint, NewIdentityServer(driver), NewControllerServer(driver), NewNodeServer(driver))
	server.Wait()
}
