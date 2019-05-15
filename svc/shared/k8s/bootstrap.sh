#!/bin/bash
set -euxo pipefail


###############################
#
#   E N A B L E   A D D O N S
#
###############################
##minikube addons enable ingress
#minikube addons enable heapster
##minikube addons enable efk
#minikube addons enable metrics-server

###############################
#
#   I N S T A L L   N A T S
#
###############################
##kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.5/00-prereqs.yaml
#kubectl apply -f nats/00-prereqs.yaml

##kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.5/10-deployment.yaml
#kubectl apply -f nats/10-deployment.yaml

#sleep 3
#kubectl create -f nats/nats_cluster.yaml


#################################
#
#   I N S T A L L   A E G O   C D
#
#################################
#kubectl create namespace argocd
#kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

#sleep 3

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

