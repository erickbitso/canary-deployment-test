#!/usr/bin/env bash

# Verify if GH_USER and GH_PAT are set, if not, exit
if [ -z "$GH_USER" ] || [ -z "$GH_PAT" ]; then
  echo "GH_USER and GH_PAT must be set"
  exit 1
fi

# Installs Istio 1.20
echo -e "\033[35mInstalling istio\033[0m"
asdf install istio 1.20.4
istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled
# Installs Argo CD
echo -e "\033[35mInstalling ArgoCD\033[0m"
kubectl create namespace argocd
kubectl apply -n argocd -f installations/argocd.yaml
# Installs Argo Rollouts
echo -e "\033[35mInstalling Argo Rollouts\033[0m"
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
# Installs Prometheus
echo -e "\033[35mInstalling Prometheus\033[0m"
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/addons/prometheus.yaml
# Installs Kiali
echo -e "\033[35mInstalling Kiali\033[0m"
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/addons/kiali.yaml
# Installs Metrics server
echo -e "\033[35mInstalling Metrics Server\033[0m"
kubectl apply -f installations/metricsserver.yaml
# Installs Keda.sh
echo -e "\033[35mInstalling Keda\033[0m"
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.14.0/keda-2.14.0-core.yaml
```

## Build services to be deployed
```
echo -e "\033[35mBuilding services\033[0m"
pack build service-a --path ./services/service-a --buildpack paketo-buildpacks/go --builder paketobuildpacks/builder-jammy-base
docker tag service-a:latest service-a:blue
docker tag service-a:latest service-a:green
kind load docker-image service-a:blue
kind load docker-image service-a:green
pack build service-b --path ./services/service-b --buildpack paketo-buildpacks/go --builder paketobuildpacks/builder-jammy-base
docker tag service-b:latest service-b:blue
docker tag service-b:latest service-b:green
kind load docker-image service-b:blue
kind load docker-image service-b:green
pack build service-b-errors --path ./services/service-b-errors --buildpack paketo-buildpacks/go --builder paketobuildpacks/builder-jammy-base
docker tag service-b-errors:latest service-b:red
kind load docker-image service-b:red


## Wait for every pod to be ready
echo -e "\033[35mWaiting for all pods to be ready\033[0m"
kubectl wait --for=condition=Ready pods --all --all-namespaces --timeout=300s


# Port forward ArgoCD
echo -e "\033[35mPort forwarding ArgoCD\033[0m"
kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 &
sleep 5
# Login to ArgoCD
echo -e "\033[35mLogging in to ArgoCD\033[0m"
admin_password=$(kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath="{.data.password}" | base64 -d)
argocd login localhost:8080 --username admin --password $admin_password --insecure
# Change admin password to "password"
echo -e "\033[35mChanging admin password to 'password'\033[0m"
argocd account update-password --current-password $admin_password --new-password password
# Add the ArgoCD repo
echo -e "\033[35mAdding ArgoCD repo\033[0m"
argocd repo add https://github.com/erickbitso/canary-deployment-test --insecure
echo
echo -e "\033[35mOpening browser to ArgoCD UI. User: admin, Password: password\033[0m"
open http://localhost:8080

echo -e "\033[35mOpening Kiali dashboard\033[0m"
istioctl dashboard kiali & 
