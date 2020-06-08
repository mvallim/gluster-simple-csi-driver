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

	drv := &Driver{
		Config: config,
	}

	glog.V(1).Infof("GlusterFS Simple CSI driver initialized")

	return drv, nil

}

// NewControllerServer
func NewControllerServer(dr *Driver) *ControllerServer {
	return &ControllerServer{
		Driver: dr,
	}
}

// NewNodeServer
func NewNodeServer(dr *Driver) *NodeServer {
	return &NodeServer{
		Driver: dr,
	}
}

// NewIdentityServer
func NewIdentityServer(dr *Driver) *IdentityServer {
	return &IdentityServer{
		Driver: dr,
	}
}

// Run
func (dr *Driver) Run() {
	srv := csicommon.NewNonBlockingGRPCServer()
	srv.Start(dr.Endpoint, NewIdentityServer(dr), NewControllerServer(dr), NewNodeServer(dr))
	srv.Wait()
}
