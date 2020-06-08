package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/spf13/cobra"

	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs"
	"github.com/mvallim/gluster-simple-csi-driver/pkg/glusterfs/config"
)

func main() {

	_ = flag.CommandLine.Parse([]string{})

	var config = config.NewConfig()

	cmd := &cobra.Command{
		Use:   "gluster-simple-csi-plugin",
		Short: "GlusterFS Simple CSI plugin",
		Run: func(cmd *cobra.Command, args []string) {

			if config.Endpoint == "" {
				config.Endpoint = os.Getenv("CSI_ENDPOINT")
			}

			if config.NodeID == "" {
				config.NodeID = os.Getenv("NODE_ID")
			}

			drv, err := glusterfs.NewDriver(config)

			if err != nil {
				glog.Fatalln(err)
			}

			drv.Run()

		},
	}

	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	cmd.PersistentFlags().StringVar(&config.NodeID, "nodeid", "", "CSI node id")
	cmd.PersistentFlags().StringVar(&config.Endpoint, "endpoint", "", "CSI endpoint")

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

}
