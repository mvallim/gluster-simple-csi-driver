package glusterfs

import (
	"context"

	"github.com/golang/glog"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type IdentityServer struct {
	*Driver
}

// GetPluginInfo
func (is *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {

	resp := &csi.GetPluginInfoResponse{
		Name:          glusterfsCSIDriverName,
		VendorVersion: glusterfsCSIDriverVersion,
	}

	glog.V(1).Infof("plugininfo response: %+v", resp)

	return resp, nil
}

// GetPluginCapabilities
func (is *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {

	resp := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}

	glog.V(1).Infof("plugin capability response: %+v", resp)

	return resp, nil
}

// Probe
func (is *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}
