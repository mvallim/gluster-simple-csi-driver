package glusterfs

import (
	"context"

	"github.com/golang/glog"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

// Common allocation units
const (
	KB int64 = 1024
	MB int64 = 1024 * KB
	GB int64 = 1024 * MB
	TB int64 = 1024 * GB

	minReplicaCount = 1
	maxReplicaCount = 10
)

// ControllerServer struct of GlusterFS CSI driver with supported methods of CSI controller server spec.
type ControllerServer struct {
	*Driver
}

// RoundUpSize calculates how many allocation units are needed to accommodate a volume of given size.
// RoundUpSize(1500 * 1000*1000, 1000*1000*1000) returns '2'
// (2 GB is the smallest allocatable volume that can hold 1500MiB)
func RoundUpSize(volumeSizeBytes int64, allocationUnitBytes int64) int64 {
	return (volumeSizeBytes + allocationUnitBytes - 1) / allocationUnitBytes
}

// RoundUpToGB rounds up given quantity upto chunks of GB
func RoundUpToGB(sizeBytes int64) int64 {
	return RoundUpSize(sizeBytes, GB)
}

// CreateVolume creates and starts the volume
func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	glog.V(2).Infof("received controller create volume request %+v", protosanitizer.StripSecrets(req))

	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be empty")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is a required field")
	}

	volumeCapabilities := req.GetVolumeCapabilities()

	if volumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "volume capabilities is a required field")
	}

	for _, cap := range volumeCapabilities {
		if cap.GetBlock() != nil {
			return nil, status.Error(codes.Unimplemented, "block volume not supported")
		}
	}

	// mountPoint := ":" + cs.Config.BlockHostPath + "/" + req.Name

	// os.MkdirAll(mountPoint, 0755)

	// servers := strings.Join(cs.Config.Servers, mountPoint+" ")

	// replicas := len(cs.Config.Servers)

	// commandCreateVolume := exec.Command("gluster", "volume", "create", req.Name, "replica", string(replicas), "arbiter 1", "transport", "tcp", servers)
	// commandStartVolume := exec.Command("gluster", "volume", "start", req.Name)

	// commandCreateVolume.Run()
	// commandStartVolume.Run()

	volSizeBytes := 1 * GB

	if capRange := req.GetCapacityRange(); capRange != nil {
		volSizeBytes = RoundUpToGB(capRange.GetRequiredBytes())
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      req.Name,
			CapacityBytes: volSizeBytes,
			VolumeContext: map[string]string{
				"glustervol":    req.Name,
				"glusterserver": cs.Servers[0],
			},
		},
	}, nil
}

// DeleteVolume deletes the given volume
func (cs *ControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return &csi.DeleteVolumeResponse{}, nil
}

// ControllerPublishVolume return Unimplemented error
func (cs *ControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerUnpublishVolume return Unimplemented error
func (cs *ControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ValidateVolumeCapabilities checks whether the volume capabilities requested are supported.
func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	glog.V(2).Infof("received controller validate volume capability request %+v", protosanitizer.StripSecrets(req))

	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is nil")
	}

	volumeID := req.GetVolumeId()

	if volumeID == "" {
		return nil, status.Error(codes.InvalidArgument, "volumeId is nil")
	}

	volumeCapabilities := req.GetVolumeCapabilities()

	if volumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "volumeCapabilities is nil")
	}

	var volumeCapabilityAccessModes []*csi.VolumeCapability_AccessMode

	for _, mode := range []csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	} {
		volumeCapabilityAccessModes = append(volumeCapabilityAccessModes, &csi.VolumeCapability_AccessMode{Mode: mode})
	}

	capabilitySupport := false

	for _, caapability := range volumeCapabilities {
		for _, volumeCapabilityAccessMode := range volumeCapabilityAccessModes {
			if volumeCapabilityAccessMode.Mode == caapability.AccessMode.Mode {
				capabilitySupport = true
			}
		}
	}

	if !capabilitySupport {
		return nil, status.Errorf(codes.NotFound, "%v not supported", req.GetVolumeCapabilities())
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: req.VolumeCapabilities,
		},
	}, nil
}

// ListVolumes returns a list of volumes
func (cs *ControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetCapacity returns the capacity of the storage pool
func (cs *ControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetCapabilities returns the capabilities of the controller service.
func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {

	functionControllerServerCapabilities := func(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
		return &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		}
	}

	var controllerServerCapabilities []*csi.ControllerServiceCapability

	for _, capability := range []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
		csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
	} {
		controllerServerCapabilities = append(controllerServerCapabilities, functionControllerServerCapabilities(capability))
	}

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: controllerServerCapabilities,
	}, nil

}

// CreateSnapshot create snapshot of an existing PV
func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// DeleteSnapshot delete provided snapshot of a PV
func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListSnapshots list the snapshots of a PV
func (cs *ControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerExpandVolume
func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetVolume
func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
