package development

import lib "github.com/archisman-mridha/x-clone/infrastructure/kubernetes/manifests/lib"

clusterName: "development"

{
	argoCD: lib.#ArgoCD @app( )

  microservices: lib.#Microservices
}
