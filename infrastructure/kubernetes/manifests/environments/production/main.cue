package production

import argocdLib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib/argocd"

{
	argoCD: argocdLib.#ArgoCD & {
		helmInstallation: values: {
			// Enable auto-scaling for the server component.
			server: autoscaling: enabled: true
		}
	} @app()
}
