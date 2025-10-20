package lib

#ArgoCD: {
  namespace: "argocd"

  helmInstallation: #HelmInstallation & {
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

  defaultProject: appProjectV1Alpha1.#AppProject & {
    metadata: {
      name: "default"
      namespace: #ArgoCD.namespace
    }

    spec: {
      sourceRepos: ["*"]
      destinations: [{ namespace: "*", server: "*" }]
      clusterResourceWhitelist: [{ group: "*", kind: "*" }]

      // Orphaned Kubernetes resource is a top-level namespaced resource which does not belong to
      // any Argo CD Application. The Orphaned Resources Monitoring feature allows detecting
      // orphaned resources, inspect/remove resources using Argo CD UI and generate a warning.
      orphanedResources: {
        warn: true
      }
    }
  }
}

let serviceMonitorEnabled = metrics: serviceMonitor: enabled: true
