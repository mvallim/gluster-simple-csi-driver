package glusterfs

import (
	"context/pkg/glusterfs/config"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	glusterfsCSIDriverName    = "org.gluster.glusterfs"
	glusterfsCSIDriverVersion = "1.0.0"
)

type Driver struct {
	*config.Config
}

func New(config *config.Config) *Driver {
	return nil
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
