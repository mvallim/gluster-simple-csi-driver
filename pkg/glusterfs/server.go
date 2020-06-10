package glusterfs

import (
	"github.com/golang/glog"

	"k8s.io/utils/mount"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

// NewDriver create new instance driver
func NewDriver(config *config.Config) (*Driver, error) {

	if config == nil {
		glog.Errorf("GlusterFS Simple CSI driver initialization failed")
		return nil, nil
	}

	driver := &Driver{
		Config: config,
	}

	glog.V(1).Infof("GlusterFS Simple CSI driver initialized")

	return driver, nil

}

// NewControllerServer create new instance controller
func NewControllerServer(driver *Driver) *ControllerServer {
	return &ControllerServer{
		Driver: driver,
	}
}

// NewNodeServer create new instance node
func NewNodeServer(driver *Driver) *NodeServer {
	return &NodeServer{
		driver:  driver,
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
