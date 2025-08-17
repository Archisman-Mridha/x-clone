package lib

#AWS: {
  namespace: "aws"

  // The AWS Cloud Controller Manager is the controller that is primarily responsible for creating
  // and updating AWS loadbalancers (classic and NLB) and node lifecycle management.
  // The controller loops that are migrating out of the kube controller manager include the route
  // controller, the service controller, the node controller, and the node lifecycle controller.
  ccm: helmInstallation: #HelmInstallation & {
    repoURL: "https://kubernetes.github.io/cloud-provider-aws"
    chart: "aws-cloud-controller-manager"
    version: "0.0.9"

    releaseName: "ccm"
    namespace: #AWS.namespace
  }

  csiDrivers: {
    // The Amazon Elastic Block Store Container Storage Interface (CSI) Driver provides a CSI
    // used by Container Orchestrators to manage the lifecycle of Amazon EBS volumes.
    ebs: helmInstallation: #HelmInstallation & {
      repoURL: "https://kubernetes-sigs.github.io/aws-ebs-csi-driver"
      chart: "aws-ebs-csi-driver"
      version: "2.47.0"

      releaseName: "ebs-csi-driver"
      namespace: #AWS.namespace

      // If the driver is able to access IMDS, it will utilize that as a preferred source of
      // metadata. The EBS CSI Driver supports IMDSv1 and IMDSv2 (and will prefer IMDSv2 if both
      // are available).
      // However, by default, IMDSv2 uses a hop limit of 1. That will prevent the driver from
      // accessing IMDSv2 if run inside a container with the default IMDSv2 configuration. So, in
      // order for the driver to access IMDS, it either must be run in host networking mode, or
      // with a hop limit of at least 2.
    }
  }
}
