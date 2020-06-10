package glusterfs

import (
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

const (
	glusterfsCSIDriverName    = "org.gluster.glusterfs"
	glusterfsCSIDriverVersion = "1.0.0"
)

// Driver is the struct embedding information about the connection to gluster cluster and configuration of CSI driver.
type Driver struct {
	*config.Config
}
