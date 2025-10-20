package production

import (
  authorizationV1 "k8s.io/api/authorization/v1"

  lib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"
)

// Cluster Autoscaler is a tool that automatically adjusts the size of the Kubernetes cluster when
// one of the following conditions is true:
//
//  (1) there are pods that failed to run in the cluster due to insufficient resources.
//
//  (2) there are nodes in the cluster that have been underutilized for an extended period of time
//      and their pods can be placed on other existing nodes.
#ClusterAutoscaler: {
  namespace: "cluster-autoscaler"

  helmInstallation: lib.#HelmInstallation & {
    repoURL: "https://kubernetes.github.io/autoscaler"
    chart: "cluster-autoscaler"
    version: "9.50.1"

    releaseName: "cluster-autoscaler"
    namespace: #ClusterAutoscaler.namespace

    values: {
      // The cluster autoscaler on Cluster API uses the cluster-api project to manage the
      // provisioning and de-provisioning of nodes within a Kubernetes cluster.
      cloudProvider: "clusterapi"
      
      // This means, the cluster in which Cluster AutoScaler will be running, is both the
      // management and the workload cluster.
      // Or, in other words, the workload cluster is managing itself using ClusterAPI.
      clusterAPIMode: "incluster-incluster"

      // Let Cluster AutoScaler auto-discover the ClusterAPI MachineDeployments.
      autoDiscovery: {
        clusterName: "x-clone-production"
        namespace: lib.#ClusterAPI.namespace
      }

      extraArgs: {
        "clusterapi-cloud-config-authoritative": true
        "cordon-node-before-terminating": true
      }
    }
  }

  // To use the opt-in support for scaling from zero as defined by the Cluster API infrastructure
  // provider, you will need to add the infrastructure machine template types to your role
  // permissions for the service account associated with the cluster autoscaler deployment.
  // The service account will need permission to get, list, and watch the infrastructure machine
  // templates for your infrastructure provider.

  clusterRole: authorizationV1.#ClusterRole & {
    apiVersion: "rbac.authorization.k8s.io/v1"
    kind: "ClusterRole"
    metadata: name: "cluster-autoscaler-capi-extension"

    rules: [
      {
        apiGroups: ["infrastructure.cluster.x-k8s.io"]
        resources: ["awsmachinetemplates"]
        verbs:     ["get", "list", "watch"]
      }
    ]
  }

  clusterRoleBinding: authorizationV1.#ClusterRole & {
    apiVersion: "rbac.authorization.k8s.io/v1"
    kind: "ClusterRoleBinding"
    metadata: name: "cluster-autoscaler-capi-extension-cluster-autoscaler"
    
    roleRef: {
      apiGroup: "rbac.authorization.k8s.io"
      kind: "ClusterRole"
      name: "cluster-autoscaler-capi-extension"
    }

    subjects: [
      {
        kind: ServiceAccount
        name: "cluster-autoscaler"
        namespace: #ClusterAutoScaler.namespace
      }
    ]
  }
}
