package prodcution

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

// OpenCost is a vendor-neutral open source project for measuring and allocating cloud
// infrastructure and container costs. It’s built for Kubernetes cost monitoring to power real-time
// cost monitoring, showback, and chargeback.
//
// OpenCost will automatically read the node information node.spec.providerID to determine the
// cloud service provider (CSP) in use. If it detects the CSP is AWS, it will attempt to pull the
// AWS on-demand pricing from the configured public API URL with no further configuration required.
#OpenCost: helmInstallation: kue.#HelmInstallation & {
  repoURL: "https://opencost.github.io/opencost-helm-chart"
  chart: "opencost"
  version: "2.2.0"

  releaseName: "opencost"
  namespace: "opencost"
}
