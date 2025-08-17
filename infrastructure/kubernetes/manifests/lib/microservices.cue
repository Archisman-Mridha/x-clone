package lib

import (
  coreV1 "k8s.io/api/core/v1"
  appsV1 "k8s.io/api/apps/v1"
)

#Microservices: {
  namespace: "microservices"

  usersMicroservice: #Microservice & {
    name: "users-microservice"
  } @app( )

  profilesMicroservice: #Microservice & {
    name: "profiles-microservice"
  } @app( )

  followshipsMicroservice: #Microservice & {
    name: "followships-microservice"
  } @app( )

  postsMicroservice: #Microservice & {
    name: "posts-microservice"
  } @app( )

  feedsMicroservice: #Microservice & {
    name: "feeds-microservice"
  } @app( )
}

#Microservice: {
  name: string
  let microserviceName = name

  version: string | *"latest"
  let microserviceVersion = version

  deployment: appsV1.#Deployment & {
    apiVersion: "apps/v1"
    kind: "Deployment"

    metadata: {
      name: microserviceName
      namespace: #Microservices.namespace
      labels: {
        "app.kubernetes.io/name": microserviceName
        "app.kubernetes.io/component": "application"
      }
    }

    spec: {
      replicas: 1

      selector: matchLabels: app: microserviceName

      template: {
        metadata: labels: app: microserviceName

        spec: {
          containers: [{
            name: microserviceName
            image: "ghcr.io/archisman-mridha/\(microserviceName):\(microserviceVersion)"

            resources: {
              requests: {
                cpu: "200m"
                memory: "256Mi"
              }
              limits: memory: "256Mi"
            }

            ports: [{
              containerPort: 4000
            }]

            volumeMounts: [{
              name: "config"
              mountPath: "/var/\(microserviceName)/config.yaml"
              subPath: "config.yaml"
              readOnly: true
            }]

            livenessProbe:  grpc: port: 4000
            readinessProbe: grpc: port: 4000
          }]

          volumes: [{
            name: "config"
            secret: {
              secretName: "\(microserviceName)-config"
            }
          }]
        }
      }
    }
  }

  service: coreV1.#Service & {
    apiVersion: "v1"
    kind: "Service"

    metadata: {
      name: microserviceName
      namespace: #Microservices.namespace
    }

    spec: {
      selector: "app.kubernetes.io/name": microserviceName

      ports: [
        {
          protocol: "TCP"
          port: "4000"
          targetPort: "4000"
        }
      ]
    }
  }
}
