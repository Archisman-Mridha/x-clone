package lib

#Cilium: {
  helmInstallation: #HelmInstallation & {
    repoURL: "https://helm.cilium.io/"
    chart: "cilium"
    version: "1.18.1"

    releaseName: "cilium"
    namespace: "cilium"

    values: {
      // Enable Cilium’s kube-proxy replacement.
      kubeProxyReplacement: "true"
      k8sServiceHost: string
      k8sServicePort: 6443

      // Enable Cilium's Ingress Controller.
      ingressController: {
        enabled: true
        default: true

        // The Ingress controller will use a shared loadbalancer for all Ingress resources.
        loadbalancerMode: "shared"

        enforceHttps: true
      }

      loadBalancer: {
        // Cilium’s eBPF kube-proxy replacement supports consistent hashing by implementing a
        // variant of The Maglev hashing in its load balancer for backend selection. This improves
        // resiliency in case of failures.
        //
        // NOTE : Maglev hashing is applied only to external (N-S) traffic. For in-cluster service
        //        connections (E-W), sockets are assigned to service backends directly, e.g. at TCP
        //        connect time, without any intermediate hop and thus are not subject to Maglev.
        algorithm: "maglev"
      }

      // Enable Gateway API support.
      gatewayAPI: enabled: true

      // One of the biggest differences between Cilium’s Ingress and Gateway API support and other
      // Ingress controllers is how closely tied the implementation is to the CNI. For Cilium,
      // Ingress and Gateway API are part of the networking stack, and so behave in a different way
      // to other Ingress or Gateway API controllers.
      //
      // Other Ingress or Gateway API controllers are generally installed as a Deployment or
      // Daemonset in the cluster, and exposed via a Loadbalancer Service or similar.
      //
      // Cilium’s Ingress and Gateway API config is exposed with a Loadbalancer or NodePort
      // service, or optionally can be exposed on the Host network also. But in all of these cases,
      // when traffic arrives at the Service’s port, eBPF code intercepts the traffic and
      // transparently forwards it to Envoy (using the TPROXY kernel facility).

      // Enable mTLS.
      // Cilium under the hood will use SPIRE (an implementation of SPIFFE) for identity
      // management.
      authentication: mutual: spire: {
        enabled: true
        install: enabled: true
      }

      // Hubble is a fully distributed networking and security observability platform, built on top
      // of Cilium and eBPF.
      hubble: {
        // By default, Hubble API operates within the scope of the individual node on which the
        // Cilium agent runs. This confines the network insights to the traffic observed by the
        // local Cilium agent.
        // Upon deploying Hubble Relay, network visibility is provided for the entire cluster or
        // even multiple clusters in a ClusterMesh scenario.
        // In this mode, Hubble data can be accessed via Hubble UI.
        relay: enabled: true

        // Hubble UI is a web interface which enables automatic discovery of the services
        // dependency graph at the L3/L4 and even L7 layer, allowing user-friendly visualization
        // and filtering of data flows as a service map.
        ui: enabled: true

        // Hubble Exporter is a feature of cilium-agent that lets you write Hubble flows to a file
        // for later consumption as logs. It supports file rotation, size limits, filters, and
        // field masks.
        export: {
          // Standard hubble exporter configuration accepts only one set of filters and requires
          // cilium pod restart to change config.
          // Dynamic flow logs allow configuring multiple filters at the same time and saving
          // output in separate files. Additionally it does not require cilium pod restarts to
          // apply changed configuration.
          dynamic: {
            enabled: true
            config: {
              enabled: true
              content: [
                {
                  name: "dropped-packet-flows"
                  filePath: "/var/run/cilium/hubble/dropped-packet-flows.log"

                  // You can view all available fields for a Hubble flow here :
                  // https://docs.cilium.io/en/stable/_api/v1/flow/README/#flowfilter.
                  includeFilters: [
                    {
                      event_type: [
                        { type: 1 } // Packet dropped.
                      ]
                    }
                  ]
                  fieldMasks: [
                    "source.namespace",
                    "source.pod_name",
                    "destination.namespace",
                    "destination.pod_name",
                    "drop_reason_desc"
                  ]
                }
              ]
            }
          }
        }
      }

      // Make Cilium components expose metrics.

                prometheus: enabled: true
      operator: prometheus: enabled: true

      hubble: metrics: enabled: true
    }
  }
}
