docker_build("users-microservice",
             context    = ".",
             dockerfile = "./backend/microservices/users/Dockerfile")
k8s_yaml("./infrastructure/kubernetes/manifests/renderred/development/resources/usersMicroservice/Deployment/microservices/users-microservice.yaml")

docker_build("profiles-microservice",
             context    = ".",
             dockerfile = "./backend/microservices/profiles/Dockerfile")
k8s_yaml("./infrastructure/kubernetes/manifests/renderred/development/resources/profilesMicroservice/Deployment/microservices/profiles-microservice.yaml")

docker_build("followships-microservice",
             context    = ".",
             dockerfile = "./backend/microservices/followships/Dockerfile")
k8s_yaml("./infrastructure/kubernetes/manifests/renderred/development/resources/followshipsMicroservice/Deployment/microservices/followships-microservice.yaml")

docker_build("posts-microservice",
             context    = ".",
             dockerfile = "./backend/microservices/posts/Dockerfile")
k8s_yaml("./infrastructure/kubernetes/manifests/renderred/development/resources/postsMicroservice/Deployment/microservices/posts-microservice.yaml")

docker_build("feeds-microservice",
             context    = ".",
             dockerfile = "./backend/microservices/feeds/Dockerfile")
k8s_yaml("./infrastructure/kubernetes/manifests/renderred/development/resources/feedsMicroservice/Deployment/microservices/feeds-microservice.yaml")
