package glusterfs

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type IdentityServer struct {
	*Driver
}

// GetPluginInfo
func (is *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return nil, nil
}

// GetPluginCapabilities
func (is *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return nil, nil
}

// Probe
func (is *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return nil, nil
}
