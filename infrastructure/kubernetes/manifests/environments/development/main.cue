package development

import lib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

let clusterName = "development"

{
	argoCD: lib.#ArgoCD @app( )

  microservices: lib.#Microservices
}
