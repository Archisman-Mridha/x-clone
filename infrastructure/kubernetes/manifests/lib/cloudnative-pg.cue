package lib

#CloudnativePG: {
  namespace: "cloudnative-pg"

  helmInstallation: #HelmInstallation & {
    repoURL: "https://cloudnative-pg.io/charts/"
    chart: "cloudnative-pg"
    version: "0.26.0"

    releaseName: "cloudnative-pg"
    namespace: #CloudnativePG.namespace

    values: {
      monitoring: {
        podMonitorEnabled: true

        grafanaDashboard: create: true
      }
    }
  }
}
