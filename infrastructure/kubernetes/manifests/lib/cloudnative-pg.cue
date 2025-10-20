package lib

#CloudnativePG: helmInstallation: #HelmInstallation & {
  repoURL: "https://cloudnative-pg.io/charts/"
  chart: "cloudnative-pg"
  version: "0.26.0"

  releaseName: "cloudnative-pg"
  namespace: "cloudnative-pg"

  values: {
    monitoring: {
      podMonitorEnabled: true

      grafanaDashboard: create: true
    }
  }
}
