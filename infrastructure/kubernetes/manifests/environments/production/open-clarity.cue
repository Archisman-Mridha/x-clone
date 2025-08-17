package prodcution

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

// OpenClarity is an open source tool for agentless detection and management of Virtual Machine
// Software Bill Of Materials (SBOM) and security threats such as vulnerabilities, exploits,
// malware, rootkits, misconfigurations and leaked secrets.
//
// OpenClarity uses a pluggable scanning infrastructure, using several tools that can be
// enabled/disabled on an individual basis.
#OpenClarity: helmInstallation: kue.#HelmInstallation & {
  repoURL: "oci://ghcr.io/openclarity/charts"
  chart: "openclarity"
  version: "1.1.2"

  releaseName: "openclarity"
  namespace: "openclarity"

  values: {
    orchestrator: {
      provider: "kubernetes"
      serviceAccount: automountServiceAccountToken: true
    }
  }
}
