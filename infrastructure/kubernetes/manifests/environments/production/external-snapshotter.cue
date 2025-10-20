package production

import kustomizeAPITypes "sigs.k8s.io/kustomize/api/types"

// The CSI external-snapshotter is part of Kubernetes implementation of Container Storage Interface
// (CSI) and implements both the volume snapshot and the volume group snapshot feature.
#ExternalSnapshotter: kustomization: kustomizeAPITypes.#Kustomization & {
  namespace: "external-snapshotter"

  resources: [
    "https://github.com/kubernetes-csi/external-snapshotter/client/config/crd",
    "https://github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/snapshot-controller"
  ]
}
