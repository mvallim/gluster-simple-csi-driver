package main

import (
	"flag"
	"log"

	"github.com/mvallim/gluster-csi-driver/pkg/glusterfs/driver"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/org.gluster.glusterfs/csi.sock", "CSI endpoint")
	)

	flag.Parse()

	drv, err := driver.NewDriver(*endpoint)

	if err != nil {
		log.Fatalln(err)
	}

	if err := drv.Run(); err != nil {
		log.Fatalln(err)
	}

}
