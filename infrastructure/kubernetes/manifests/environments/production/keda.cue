package prodcution

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

#KEDA: {
  helmInstallation: kue.#HelmInstallation & {
    repoURL: ""
    chart: "keda"
    version: ""

    releaseName: "keda"
    namespace: "keda"

    values: { }
  }
}
