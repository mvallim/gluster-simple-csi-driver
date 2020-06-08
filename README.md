<p align="center">
  <img src="images/csi-logo.png">
</p>

# GlusterFS CSI Driver

_This repo contains CSI driver for Gluster. The Container Storage Interface (CSI) is a proposed new industry standard for cluster-wide volume plugins. “Container Storage Interface” (CSI) enables storage vendors (SP) to develop a plugin once and have it work across a number of container orchestration (CO) systems._

## Why CSI?

_Although prior to CSI Kubernetes provided a powerful volume plugin system, it was challenging to add support for new volume plugins to Kubernetes: volume plugins were “in-tree” meaning their code was part of the core Kubernetes code and shipped with the core Kubernetes binaries—vendors wanting to add support for their storage system to Kubernetes (or even fix a bug in an existing volume plugin) were forced to align with the Kubernetes release process. In addition, third-party storage code caused reliability and security issues in core Kubernetes binaries and the code was often difficult (and in some cases impossible) for Kubernetes maintainers to test and maintain._

_CSI was developed as a standard for exposing arbitrary block and file storage storage systems to containerized workloads on Container Orchestration Systems (COs) like Kubernetes. With the adoption of the Container Storage Interface, the Kubernetes volume layer becomes truly extensible. Using CSI, third-party storage providers can write and deploy plugins exposing new storage systems in Kubernetes without ever having to touch the core Kubernetes code. This gives Kubernetes users more options for storage and makes the system more secure and reliable._
