package argocd

import "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

#ArgoCD: {
  namespace: string | *"argocd"

  helmInstallation: lib.#HelmInstallation & {
    repoURL: "https://argoproj.github.io/argo-helm"
    chart: "argo-cd"
    version: "8.2.7"

    releaseName: "argocd"
    namespace: #ArgoCD.namespace

    values: {
      // Enable monitoring for all the components.
      server: serviceMonitorEnabled & {...}
      controller: serviceMonitorEnabled
      dex: serviceMonitorEnabled
      redis: serviceMonitorEnabled
      repoServer: serviceMonitorEnabled
      notifications: serviceMonitorEnabled
    }
  }

  serviceMonitorEnabled: metrics: serviceMonitor: enabled: true
}
