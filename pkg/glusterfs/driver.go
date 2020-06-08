package glusterfs

import (
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

// New
func New(config *config.Config) (*Driver, error) {
	return nil, nil
}

// NewControllerServer
func NewControllerServer(dr *Driver) *ControllerServer {
	return nil
}

// NewNodeServer
func NewNodeServer(dr *Driver) *NodeServer {
	return nil
}

// NewIdentityServer
func NewIdentityServer(dr *Driver) *IdentityServer {
	return nil
}

// Run
func (dr *Driver) Run() {
	srv := csicommon.NewNonBlockingGRPCServer()
	srv.Start("g.Endpoint", NewIdentityServer(dr), NewControllerServer(dr), NewNodeServer(dr))
	srv.Wait()
}
