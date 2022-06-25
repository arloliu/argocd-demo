#!/usr/bin/env bash
# minikube version: v1.26.0
# nginx ingress controller: v1.2.1
export ARGOCD_VERSION=v2.3.3
CUR_DIR=$(dirname $0)

if [ ! kubectl get namespace/argocd > /dev/null 2>&1 ]; then
    echo "Create argocd namespace"
    kubectl create namespace argocd
fi

if [ $(minikube addons list -o json | jq .ingress.Status) = "\"disabled\"" ]; then
    echo "enable minikube ingress addon"
    minikube addons enable ingress
fi

if [ ! kubectl get -n argocd svc/argocd-server > /dev/null 2>&1 ]; then
    echo "Install argocd $ARGOCD_VERSION"
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/${ARGOCD_VERSION}/manifests/install.yaml
fi

