package glusterfs

import (
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/mvallim/gluster-csi-driver/pkg/glusterfs/config"
)

const (
	glusterfsCSIDriverName    = "org.gluster.glusterfs"
	glusterfsCSIDriverVersion = "1.0.0"
)

type Driver struct {
	*config.Config
}

func New(config *config.Config) (*Driver, error) {
	return nil, nil
}

func NewControllerServer(dr *Driver) *ControllerServer {
	return nil
}

func NewNodeServer(dr *Driver) *NodeServer {
	return nil
}

func NewIdentityServer(dr *Driver) *IdentityServer {
	return nil
}

func (dr *Driver) Run() {
	srv := csicommon.NewNonBlockingGRPCServer()
	srv.Start("g.Endpoint", NewIdentityServer(dr), NewControllerServer(dr), NewNodeServer(dr))
	srv.Wait()
}
