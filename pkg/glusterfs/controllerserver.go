package glusterfs

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog"

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
	client     kubernetes.Interface
	restclient rest.Interface
	config     *rest.Config
}

// RoundUpSize calculates how many allocation units are needed to accommodate a volume of given size.
// RoundUpSize(1500 * 1000*1000, 1000*1000*1000) returns '2'
// (2 GB is the smallest allocatable volume that can hold 1500MiB)
func RoundUpSize(volumeSizeBytes int64, allocationUnitBytes int64) int64 {
	return (volumeSizeBytes + allocationUnitBytes - 1) / allocationUnitBytes
}

// RoundUpToGB rounds up given quantity upto chunks of GB
func RoundUpToGB(sizeBytes int64) int64 {
	return RoundUpSize(sizeBytes, GB) * GB
}

func (cs *ControllerServer) selectPod(host string) (*v1.Pod, error) {

	podList, err := cs.client.CoreV1().Pods("glusterfs").List(meta_v1.ListOptions{LabelSelector: cs.ServerLabel})

	if err != nil {
		return nil, err
	}

	pods := podList.Items

	if len(pods) == 0 {
		return nil, fmt.Errorf("No pods found for glusterfs, LabelSelector: %v", cs.ServerLabel)
	}

	for _, pod := range pods {
		if pod.Status.PodIP == host {
			klog.Infof("Pod selecterd: %v/%v\n", pod.Namespace, pod.Name)
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("No pod found to match NodeName == %s", host)
}

func (cs *ControllerServer) ExecuteCommand(pod *v1.Pod, commands []string) error {

	for _, command := range commands {

		klog.Infof("Pod: %s, ExecuteCommand: %s", pod.Name, command)

		containerName := pod.Spec.Containers[0].Name

		req := cs.restclient.Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(pod.Namespace).
			SubResource("exec").
			Param("container", containerName).
			Param("stdout", "true").
			Param("stderr", "true").
			Param("command", "/bin/bash").
			Param("command", "-c").
			Param("command", command)

		exec, err := remotecommand.NewSPDYExecutor(cs.config, "POST", req.URL())

		if err != nil {
			klog.Fatalf("Failed to create NewExecutor: %v", err)
			return err
		}

		var b bytes.Buffer
		var berr bytes.Buffer

		err = exec.Stream(remotecommand.StreamOptions{
			Stdout: &b,
			Stderr: &berr,
			Tty:    false,
		})

		klog.Infof("Result: %v", b.String())
		klog.Infof("Result: %v", berr.String())

		if err != nil {
			klog.Errorf("Failed to create Stream: %v", err)
			return err
		}
	}

	return nil
}

// CreateVolume creates and starts the volume
func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	klog.Infof("received controller create volume request %+v", protosanitizer.StripSecrets(req))

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

	pvcName := req.Parameters["csi.storage.k8s.io/pvc/name"]
	pvcNameSpace := req.Parameters["csi.storage.k8s.io/pvc/namespace"]
	pvName := req.Name

	hosts := strings.Split(cs.Servers, ";")
	path := filepath.Join(cs.HostPath, pvcNameSpace, pvcName, pvName)

	var bricks string

	for _, host := range hosts {
		bricks += strings.Join([]string{host, path + " "}, ":")
	}

	pod, err := cs.selectPod(hosts[0])

	if err != nil {
		return nil, err
	}

	volSizeBytes := 1 * GB

	if capRange := req.GetCapacityRange(); capRange != nil {
		volSizeBytes = RoundUpToGB(capRange.GetRequiredBytes())
	}

	commands := []string{
		fmt.Sprintf("gluster --mode=script volume create %s replica %v arbiter 1 transport tcp %s", req.Name, len(hosts), bricks),
		fmt.Sprintf("gluster --mode=script volume start %s", req.Name),
	}

	err = cs.ExecuteCommand(pod, commands)

	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      req.Name,
			CapacityBytes: int64(volSizeBytes),
			VolumeContext: map[string]string{
				"glustervol":    req.Name,
				"glusterserver": hosts[0],
				"glusterpath":   path,
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

	klog.Infof("received controller validate volume capability request %+v", protosanitizer.StripSecrets(req))

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
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
		csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
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

// ControllerExpandVolume resizes a volume
func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetVolume
func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
