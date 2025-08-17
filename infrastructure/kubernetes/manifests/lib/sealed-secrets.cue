package lib

// Sealed Secrets is used to encrypt a Kubernetes Secret into a SealedSecret, which you can store
// in a public Git repository. The SealedSecret can be decrypted only by the Sealed Secrets
// controller (which has access to the encryption key), running in the target cluster.
#SealedSecrets: {
  namespace: "sealed-secrets"

  helmInstallation: #HelmInstallation & {
    repoURL: "https://bitnami-labs.github.io/sealed-secrets"
    chart: "sealed-secrets"
    version: "2.17.3"

    releaseName: "sealed-secrets"
    namespace: #SealedSecrets.namespace

    values: {
      namespace: #SealedSecrets.namespace

      // By default, the kubeseal command line interface (CLI) tries to access the controller with
      // the name sealed-secrets-controller. 
      fullnameOverride: "sealed-secrets-controller"
    }
  }
}
