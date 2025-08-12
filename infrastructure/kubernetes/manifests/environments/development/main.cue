package production

import argocdLib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib/argocd"

{
	argoCD: argocdLib.#ArgoCD @app()
}
