package lib

// ExternalDNS allows you to control DNS records dynamically via Kubernetes resources in a DNS
// provider-agnostic way.
#ExternalDNS: {
  clusterName: string

  namespace: "external-dns"

  helmInstallation: #HelmInstallation & {
    repoURL: "https://kubernetes-sigs.github.io/external-dns/"
    chart: "external-dns"
    version: "1.17.0"

    releaseName: "external-dns"
    namespace: #ExternalDNS.namespace

    values: {
      sources: [
        "gateway-tcproute"
        "gateway-httproute"
        "gateway-grpcroute"
        "gateway-udproute"
      ]

      provider: name: "cloudflare"

      env: [
        {
          name: "CF_API_TOKEN"
          valueFrom: secretKeyRef: {
            name: "cloudflare-api-credentials"
            key: "api-key"
          }
        }
      ]

      domainFilters: ["\(#ExternalDNS.clusterName).projectsofarchi.xyz"]
      txtOwnerId: #ExternalDNS.clusterName

      // Enable exposing metrics.
      serviceMonitor: enabled: true
    }
  }
}
