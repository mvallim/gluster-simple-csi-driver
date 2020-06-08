package main

import (
	"log"

	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs"
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

func main() {

	var config = config.NewConfig()

	config.NodeID = ""
	config.Endpoint = "unix:///var/lib/kubelet/plugins/org.gluster.glusterfs/csi.sock"

	drv, err := glusterfs.New(config)

	if err != nil {
		log.Fatalln(err)
	}

	drv.Run()

}
