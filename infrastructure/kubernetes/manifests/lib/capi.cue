package lib

import infrastructureProviderV1Alpha2 "operator.cluster.x-k8s.io/infrastructureprovider/v1alpha2"

// Cluster API is a Kubernetes subproject focused on providing declarative APIs and tooling to
// simplify provisioning, upgrading, and operating multiple Kubernetes clusters.
#ClusterAPI: {
  namespace: "cluster-api"

  // The Cluster API Operator is a Kubernetes Operator designed to empower cluster administrators
  // to handle the lifecycle of Cluster API providers within a management cluster using a
  // declarative approach.
  operator: helmInstallation: #HelmInstallation & {
    repoURL: "https://kubernetes-sigs.github.io/cluster-api-operator"
    chart: "cluster-api-operator"
    version: "0.22.0"

    releaseName: "cluster-api-operator"
    namespace: #ClusterAPI.namespace

    values: {
      core: "cluster-api": {
        namespace: #ClusterAPI.namespace
        version: "v1.10.4"
      }

      bootstrap: kubeadm: {
        namespace: #ClusterAPI.namespace
        version: "v1.10.4"
      }

      controlPlane: kubeadm: {
        namespace: #ClusterAPI.namespace
        version: "v1.10.4"
      }
    }
  }

  infrastructureProvider: infrastructureProviderV1Alpha2.#InfrastructureProvider & {
    let capaVersion = "v2.8.4"

    metadata: {
      name: "aws"
      namespace: #ClusterAPI.namespace
    }

    spec: {
      version: capaVersion

      fetchConfig: url: "https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/\(capaVersion)/infrastructure-components.yaml"

      configSecret: {
        name: "aws-credentials"
        namespace: #ClusterAPI.namespace
      }

      manager: {
        // Restricts the manager's cache to watch objects in the desired namespace.
        cacheNamespace: #ClusterAPI.namespace
      }
    }
  }
}
