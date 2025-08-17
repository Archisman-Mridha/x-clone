package prodcution

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

// Harbor is an open source registry that secures artifacts with policies and role-based access
// control, ensures images are scanned and free from vulnerabilities, and signs images as trusted. 
#Harbor: {
  namespace: "harbor"

  helmInstallation: kue.#HelmInstallation & {
    repoURL: "https://helm.goharbor.io"
    chart: "harbor"
    version: "1.17.2"

    releaseName: "harbor"
    namespace: #Harbor.namespace

    values: { }
  }
}
