#!/usr/bin/env bash
# minikube version: v1.26.0
# nginx ingress controller: v1.2.1
export ARGOCD_VERSION=v2.3.3


kubectl create namespace argocd

kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/${ARGOCD_VERSION}/manifests/install.yaml

# nginx ingress controller
minikube addons enable ingress

# setup ingress for argocd.local
kubectl apply -f ./ingress.yaml

INGRESS_HOST=$(minikube ip)
echo "Ingress IP: ${INGRESS_HOST}"
