package glusterfs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/utils/mount"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"k8s.io/klog"
)

// NodeServer struct of Glusterfs CSI driver with supported methods of CSI node server spec.
type NodeServer struct {
	*Driver
	mounter mount.Interface
}

// NodeStageVolume mounts the volume to a staging path on the node.
func (ns *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

// NodeUnstageVolume unstages the volume from the staging path
func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

// NodePublishVolume mounts the volume mounted to the staging path to the target path
func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {

	klog.Infof("received node publish volume request %+v", protosanitizer.StripSecrets(req))

	targetPath := req.GetTargetPath()

	mounted, error := ns.mounter.IsLikelyNotMountPoint(targetPath)

	if error != nil {

		if os.IsNotExist(error) {
			if error := os.MkdirAll(targetPath, 0777); error != nil {
				return nil, status.Error(codes.Internal, error.Error())
			}
			mounted = true
		} else {
			return nil, status.Error(codes.Internal, error.Error())
		}
	}

	if !mounted {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	mountFlags := req.GetVolumeCapability().GetMount().GetMountFlags()

	if req.GetReadonly() {
		mountFlags = append(mountFlags, "rw", "acl")
	}

	server := req.GetVolumeContext()["glusterserver"]
	brick := req.GetVolumeContext()["glustervol"]
	source := fmt.Sprintf("%s:%s", server, brick)

	error = ns.mounter.Mount(source, targetPath, "glusterfs", mountFlags)

	if error != nil {

		if os.IsPermission(error) {
			return nil, status.Error(codes.PermissionDenied, error.Error())
		}

		if strings.Contains(error.Error(), "invalid argument") {
			return nil, status.Error(codes.InvalidArgument, error.Error())
		}

		return nil, status.Error(codes.Internal, error.Error())
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {

	klog.Infof("received node unpublish volume request %+v", protosanitizer.StripSecrets(req))

	targetPath := req.GetTargetPath()

	mounted, error := ns.mounter.IsLikelyNotMountPoint(targetPath)

	if error != nil {

		if os.IsNotExist(error) {
			return nil, status.Error(codes.NotFound, "Targetpath not found")
		}

		return nil, status.Error(codes.Internal, error.Error())
	}

	if mounted {
		return nil, status.Error(codes.NotFound, "Volume not mounted")
	}

	error = mount.CleanupMountPoint(targetPath, ns.mounter, false)

	if error != nil {
		return nil, status.Error(codes.Internal, error.Error())
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats returns volume capacity statistics available for the volume
func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// NodeExpandVolume expanding the file system on the node
func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// NodeGetCapabilities returns the supported capabilities of the node server
func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

// NodeGetInfo returns NodeGetInfoResponse for CO.
func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {

	klog.Infof("received node get info request %+v", protosanitizer.StripSecrets(req))

	return &csi.NodeGetInfoResponse{
		NodeId: ns.NodeID,
	}, nil

}
