package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
			handle(config)
		},
	}

	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	cmd.PersistentFlags().StringVar(&config.NodeID, "nodeid", "", "CSI node id")
	_ = cmd.MarkPersistentFlagRequired("nodeid")
	cmd.PersistentFlags().StringVar(&config.Endpoint, "endpoint", "", "CSI endpoint")
	_ = cmd.MarkPersistentFlagRequired("endpoint")

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

}

func handle(config *config.Config) {

	if config.Endpoint == "" {
		config.Endpoint = os.Getenv("CSI_ENDPOINT")
	}

	if config.NodeID == "" {
		config.NodeID = os.Getenv("NODE_ID")
	}

	drv, err := glusterfs.NewDriver(config)

	if err != nil {
		log.Fatalln(err)
	}

	drv.Run()

}
