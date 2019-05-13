#!/bin/bash
set -euxo pipefail

###############################
#
#   I N S T A L L   A E G O C D
#
###############################
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Initial password is the ArgoCD server pod id
TEMP_PASSWORD=$(kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2)

# Port-forward to access the Argo CD server locally on port 8080:
kubectl port-forward -n argocd svc/argocd-server 8080:443 &
sleep 3

# Login and update password
argocd login :8080 --insecure --username admin --password $TEMP_PASSWORD
argocd account update-password --current-password $TEMP_PASSWORD --new-password $ARGOCD_PASSWORD

# Add Delinkcious services
python bootstrap_argocd.py

###############################
#
#   I N S T A L L   N A T S
#
###############################
kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.5/00-prereqs.yaml
kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.5/10-deployment.yaml

kubectl create -f nats_cluster.yaml
