package lib

// Crossplane is an open source Kubernetes extension that transforms your Kubernetes cluster into a
// universal control plane. It lets you manage anything, anywhere, all through standard Kubernetes
// APIs.
//
// NOTE : A control plane creates and manages the lifecycle of resources. It constantly checks that
//        the intended resources exist, reports when the intended state doesn’t match reality and
//        acts to make things right.
#Crossplane: {
  namespace: "crossplane"

  helmInstallation: kue.#HelmInstallation & {
    repoURL: "https://charts.crossplane.io/stable"
    version: "1.10.0"
    chartPath: "crossplane"

    releaseName: "crossplane"
    namespace: #Crossplane.namespace

    values: {
      metrics: enabled: true
    }
  }

  aws: { }
}
