#!/bin/bash
set -euxo pipefail


#################################################
#
#   E N A B L E   A D D O N S ( M I N I K U B E )
#
#################################################
##
#minikube addons enable ingress
#minikube addons enable heapster
#minikube addons enable efk
#minikube addons enable metrics-server

################################
#
#   I N S T A L L   H E L M
#
################################
#kubectl apply -f helm_rbac.yaml
#helm init --service-account tiller

###############################################
#
#   I N S T A L L   M E T R I C S   S E R V E R
#
###############################################
#helm install stable/metrics-server \
#    --name metrics-server          \
#    --version 2.0.4                \
#    --namespace kube-system

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
#   I N S T A L L   A R G O   C D
#
#################################
#kubectl create namespace argocd
#kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Initial password is the ArgoCD server pod id
#ARGOCD_POD_NAME=$(kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2)
#while [ ARGOCD_POD_NAME == "" ]
#do
#  ARGOCD_POD_NAME=$(kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2)
#  sleep 1
#done

# Wait for pod to be running
#PHASE=$(kubectl -n argocd get po $ARGOCD_POD_NAME -o jsonpath='{.status.phase}')
#while [ $PHASE != "Running" ]
#do
#  PHASE=$(kubectl -n argocd get po $ARGOCD_POD_NAME -o jsonpath='{.status.phase}')
#  sleep 1
#done

# Port-forward to access the Argo CD server locally on port 8080:
#kubectl port-forward -n argocd svc/argocd-server 8080:443 &
#sleep 3

# Login and update password
#argocd login :8080 --insecure --username admin --password $ARGOCD_POD_NAME
#argocd account update-password --current-password $ARGOCD_POD_NAME --new-password $ARGOCD_PASSWORD

# Add Delinkcious services
#python bootstrap_argocd.py

#################################
#
#   I N S T A L L   N U C L I O
#
#################################
# Install nuclio in its own namespace
#kubectl create namespace nuclio

#kubectl apply -f https://raw.githubusercontent.com/nuclio/nuclio/master/hack/k8s/resources/nuclio-rbac.yaml
#kubectl apply -f nuclio/nuclio-rbac.yaml

#kubectl apply -f https://raw.githubusercontent.com/nuclio/nuclio/master/hack/k8s/resources/nuclio.yaml
#kubectl apply -f nuclio/nuclio.yaml

## Get nuctl CLI and create a symlink
#curl -LO https://github.com/nuclio/nuclio/releases/download/1.1.5/nuctl-1.1.5-darwin-amd64
#sudo mv nuctl-1.1.5-darwin-amd64 /usr/local/bin/nuctl
#chmod +x nuctl-1.1.5-darwin-amd64


# Create an image pull secret, so Nuclio can deploy functions to our cluster.
kubectl create secret docker-registry registry-credentials -n nuclio \
  --docker-username g1g1 \
  --docker-password $DOCKERHUB_PASSWORD \
  --docker-server registry.hub.docker.com \
  --docker-email the.gigi@gmail.com


# Deploy the link checker nuclio function
pushd .
cd ../../../fun/link_checker
nuctl deploy link-checker -n nuclio -p . --registry g1g1
popd


#######################################
#
#   I N S T A L L   P R O M E T H E U S
#
#######################################
helm install --name prometheus stable/prometheus
