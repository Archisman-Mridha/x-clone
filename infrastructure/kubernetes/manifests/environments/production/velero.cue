package production

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

#Velero: {
  namespace: "velero"

  helmInstallation: lib.#HelmInstallation & {
    repoURL: "https://vmware-tanzu.github.io/helm-charts/"
    chart: "velero"
    version: "10.1.0"

    releaseName: "velero"
    namespace: #Velero.namespace

    values: { }
  }
}
