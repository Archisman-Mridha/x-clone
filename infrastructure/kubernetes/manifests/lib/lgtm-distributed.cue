package lib

#LGTMDistribted: helmInstallation: #HelmInstallation & {
  repoURL: "https://grafana.github.io/helm-charts"
  chart: "lgtm-distributed"
  version: "2.1.0"

  releaseName: "lgtm-distributed"
  namespace: "grafana"

  values: {
    grafana: enabled: false
  }
}
