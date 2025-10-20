package production

import lib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

let clusterName = "production"

{
  networking: {
    cilium: lib.#Cilium & {
      helmInstallation: values: k8sServiceHost: string
    } @app( )

    certManager: lib.#CertManager & {
      cloudflareCredentialsSealedSecret: _
    } @app( )

    gatewayAPI: lib.#GatewayAPI @app( )

    externalDNS: lib.#ExternalDNS & {
      clusterName: clusterName
    } @app( )
  }

  aws: lib.#AWS @app( )

  gitOps: {
    argoCD: lib.#ArgoCD & {
      helmInstallation: values: {
        // Enable auto-scaling for the server component.
        server: autoscaling: enabled: true
      }
    } @app( )

    sealedSecrtes: lib.#SealedSecrets @app( )

    clusterAPI: lib.#ClusterAPI @app( )

    crossplane: lib.#CrossPlane @app( )
  }

  clusterAutoscaler: #ClusterAutoscaler @app( )

  harbor: #Harbor @app( )

  microservices: lib.#Microservices

  keda: #KEDA @app( )

  monitoring: {
    lgtmDistributed: lib.#LGTMDistributed & {
      helmInstallation: values: {
        grafana: {
          autoscaling: enabled: true

          // By default, persistent storage is disabled, which means that Grafana uses ephemeral
          // storage, and all data will be stored within the container’s file system. This data
          // will be lost if the container is stopped, restarted, or if the container crashes.
          persistence: {
            enabled: true
            size: "1Gi"
          }
        }
      }
    } @app( )

    openClarity: #OpenClarity @app( )

    openCost: #OpenCost @app( )
  }

  backup: {
    velero: #Velero @app( )

    externalSnapshotter: #ExternalSnapshotter @app( )
  }
}
