package glusterfs

import (
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
