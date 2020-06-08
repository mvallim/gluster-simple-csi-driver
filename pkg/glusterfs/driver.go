package glusterfs

import (
	"github.com/golang/glog"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

const (
	glusterfsCSIDriverName    = "org.gluster.glusterfs"
	glusterfsCSIDriverVersion = "1.0.0"
)

// Driver
type Driver struct {
	*config.Config
}

// NewDriver
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

// NewControllerServer
func NewControllerServer(driver *Driver) *ControllerServer {
	return &ControllerServer{
		Driver: driver,
	}
}

// NewNodeServer
func NewNodeServer(driver *Driver) *NodeServer {
	return &NodeServer{
		Driver: driver,
	}
}

// NewIdentityServer
func NewIdentityServer(driver *Driver) *IdentityServer {
	return &IdentityServer{
		Driver: driver,
	}
}

// Run
func (driver *Driver) Run() {
	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(driver.Endpoint, NewIdentityServer(driver), NewControllerServer(driver), NewNodeServer(driver))
	server.Wait()
}
