## Kubernetes native development

Bootstrap a local K3D cluster :
```shell
k3d cluster create --config k3d/config.yaml
```

Create the root ArgoCD Application :
```shell
kubectl apply -f renderred/development/applications/root.yaml
```

## REFERENCEs

- [Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)

- [Scaling in the Clouds: Istio Ambient vs. Cilium](https://istio.io/latest/blog/2024/ambient-vs-cilium/)
