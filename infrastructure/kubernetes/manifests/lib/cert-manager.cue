package lib

import (
  clusterIssuerV1 "cert-manager.io/clusterissuer/v1"
  certificateV1 "cert-manager.io/certificate/v1"

  sealedSecretV1Alpha1 "bitnami.com/sealedsecret/v1alpha1"
)

// cert-manager creates TLS certificates for workloads in your Kubernetes or OpenShift cluster and
// renews the certificates before they expire. It can obtain certificates from a variety of
// certificate authorities, including: Let's Encrypt, HashiCorp Vault, Venafi and private PKI.
#CertManager: {
  namespace: "cert-manager"

  helmInstallation: #HelmInstallation & {
    repoURL: "https://charts.jetstack.io"
    chart: "cert-manager"
    version: "1.18.2"

    releaseName: "cert-manager"
    namespace: #CertManager.namespace

    values: {
      crds: enabled: true

      // Enable exposing metrics.
      prometheus: {
        enabled: true
        serviceMonitor: enabled: true
      }
    }
  }

  let cloudflareCredentialsSecretName = "cloudflare-credentials"

  // ClusterIssuer represents a certificate authority (CA) able to sign (issue) certificates in
  // response to Certificate Signing Requests (CSRs).
  // The ClusterIssuer resource is cluster scoped. This means that when referencing a secret via
  // the secretName field, secrets will be looked for in the Cluster Resource Namespace. By
  // default, this namespace is cert-manager.
  clusterIssuer: clusterIssuerV1.#ClusterIssuer & {
    metadata: name: "letsencrypt"

    spec: {
			// The ACME Issuer type represents a single account registered with the Automated Certificate
			// Management Environment (ACME) Certificate Authority server. Certificates issued by public
      // ACME servers are typically trusted by client's computers by default.
      acme: {
        server: "https://acme-v02.api.letsencrypt.org/directory"

        // Let's Encrypt will use this to contact you about expiring certificates, and issues
        // related to your account.
        email: "archismanmridha12345@gmail.com"

        // When you create a new ACME Issuer, CertManager will generate a private key which is used
        // to identify you with the ACME server.
        privateKeySecretRef: {
          // Name of the Kubernetes Secret that will be used to store the private key.
          name: "letsencrypt-private-key"
        }

        // In order for the ACME CA server to verify that a client owns the domain(s) a
				// certificate is being requested for, the client must complete challenges.
				solvers: [{
					// DNS challenge solver will be used to create and manage DNS records. When a certificate
					// is requested, CertManager will create a DNS TXT record called '_acme-challenge' under
					// the domain you want a certificate for. The ACME CA (e.g. Let's Encrypt) will check for
					// the existence of that TXT record to verify that you control the domain. Once verified,
					// the CA will issue the TLS certificate. CertManager will then clean up the DNS records
					// it created.
					dns01: {
						cloudflare: {
							email: "archismanmridha12345@gmail.com"
							apiTokenSecretRef: {
								name: "cloudflare-credentials"
								key: "api-token"
							}
						}
					}
				}]
      }
    }
  }

  cloudflareCredentialsSealedSecret: sealedSecretV1Alpha1.#SealedSecret & {
    metadata: {
      name: "cloudflare-credentials"
      namespace: #CertManager.namespace
    }
  }

  wildcardCertificate: certificateV1.#Certificate & {
    metadata: {
      name: "wildcard"
      namespace: #GatewayAPI.namespace
    }

    spec: {
      secretName: "wildcard-certificate-tls-keys"

      dnsNames: [
        "\(#ExternalDNS.clusterName).projectsofarchi.xyz"
        "*.\(#ExternalDNS.clusterName).projectsofarchi.xyz"
      ]

      issuerRef: {
        kind: "ClusterIssuer"
        name: "letsencrypt"
      }
    }
  }
}
